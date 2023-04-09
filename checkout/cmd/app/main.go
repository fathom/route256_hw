package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	CheckoutV1 "route256/checkout/internal/api/checkout_v1"
	LomsClient "route256/checkout/internal/clients/grpc/loms_client"
	ProductClient "route256/checkout/internal/clients/grpc/product_client"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/interceptors"
	"route256/checkout/internal/logger"
	"route256/checkout/internal/metrics"
	db "route256/checkout/internal/repository/db_repository"
	"route256/checkout/internal/repository/db_repository/transactor"
	desc "route256/checkout/pkg/checkout_v1"
	"time"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// Checkout
// Сервис отвечает за корзину и оформление заказа.

func main() {

	// Инициализация конфига
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	// Инициализация логирования
	logger.Init(config.ConfigData.Dev)

	logger.Info(
		"init config services",
		zap.String("loms", config.ConfigData.Services.Loms),
		zap.String("product", config.ConfigData.Services.Product),
	)

	// Инициализация grpc сервера
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", config.ConfigData.GrpcPort))
	if err != nil {
		logger.Fatal("failed start listen", zap.Error(err))
		panic(err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptors.LoggingInterceptor(logger.GetLogger()),
				interceptors.Metrics,
			),
		),
	)

	reflection.Register(s)

	// Инициализация базы данных
	dbpool, err := pgxpool.Connect(context.Background(), config.ConfigData.DatabaseURL)
	if err != nil {
		logger.Fatal("unable to create connection pool", zap.Error(err))
		panic(err)
	}
	defer dbpool.Close()

	// Инициализация зависимостей, репозиториев, транзактора
	tm := transactor.NewTransactionManager(dbpool)
	cartRepo := db.NewCartRepository(tm)

	// Инициализация LOMS клиента
	lomsConn, err := grpc.Dial(
		config.ConfigData.Services.Loms,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors.ClientInterceptor),
	)
	if err != nil {
		logger.Fatal("failed connect to loms server", zap.Error(err))
		panic(err)
	}
	defer lomsConn.Close()

	lomsClient := LomsClient.New(lomsConn)

	// Инициализация Products клиента
	productConn, err := grpc.Dial(
		config.ConfigData.Services.Product,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptors.ClientInterceptor),
	)
	if err != nil {
		logger.Fatal("failed connect to product server", zap.Error(err))
		panic(err)
	}
	defer productConn.Close()

	productServiceClient := ProductClient.New(productConn, config.ConfigData.Token)

	// Ограничиваем кол-во запросов 10rps
	limiter := rate.NewLimiter(rate.Every(1*time.Second/10), 10)

	// Инициализация домена
	businessLogic := domain.New(
		lomsClient,
		productServiceClient,
		tm,
		cartRepo,
		limiter,
	)

	// Запуск grpc сервера
	desc.RegisterCheckoutV1Server(s, CheckoutV1.NewCheckoutV1(businessLogic))
	logger.Info("grpc server listening at port", zap.String("port", config.ConfigData.GrpcPort))

	grpcPrometheus.Register(s)

	go func() {
		http.Handle("/metrics", metrics.NewMetricHandler())
		err := http.ListenAndServe(":8082", nil)
		logger.Fatal("failed metrics", zap.Error(err))
		panic(err)
	}()

	if err = s.Serve(lis); err != nil {
		logger.Fatal("failed start serve", zap.Error(err))
		panic(err)
	}
}
