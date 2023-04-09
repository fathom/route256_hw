package main

import (
	"context"
	"fmt"
	"log"
	"net"
	CheckoutV1 "route256/checkout/internal/api/checkout_v1"
	LomsClient "route256/checkout/internal/clients/grpc/loms_client"
	ProductClient "route256/checkout/internal/clients/grpc/product_client"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/interceptors"
	lg "route256/checkout/internal/logger"
	db "route256/checkout/internal/repository/db_repository"
	"route256/checkout/internal/repository/db_repository/transactor"
	desc "route256/checkout/pkg/checkout_v1"
	"time"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
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
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	logger := lg.NewLogger(config.ConfigData.Dev)
	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Fatal("logger sync", err)
		}
	}()

	logger.Info(
		"init config services",
		zap.String("loms", config.ConfigData.Services.Loms),
		zap.String("product", config.ConfigData.Services.Product),
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", config.ConfigData.GrpcPort))
	if err != nil {
		logger.Fatal("failed start listen", zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptors.LoggingInterceptor(logger),
			),
		),
	)

	reflection.Register(s)

	dbpool, err := pgxpool.Connect(context.Background(), config.ConfigData.DatabaseURL)
	if err != nil {
		logger.Fatal("unable to create connection pool", zap.Error(err))
	}
	defer dbpool.Close()

	tm := transactor.NewTransactionManager(dbpool)
	cartRepo := db.NewCartRepository(tm)

	lomsConn, err := grpc.Dial(config.ConfigData.Services.Loms, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("failed connect to loms server", zap.Error(err))
	}
	defer lomsConn.Close()

	lomsClient := LomsClient.New(lomsConn)

	productConn, err := grpc.Dial(config.ConfigData.Services.Product, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("failed connect to product server", zap.Error(err))
	}
	defer productConn.Close()

	productServiceClient := ProductClient.New(productConn, config.ConfigData.Token)

	// Ограничиваем кол-во запросов 10rps
	limiter := rate.NewLimiter(rate.Every(1*time.Second/10), 10)

	businessLogic := domain.New(
		lomsClient,
		productServiceClient,
		tm,
		cartRepo,
		limiter,
	)

	desc.RegisterCheckoutV1Server(s, CheckoutV1.NewCheckoutV1(businessLogic))
	logger.Info("grpc server listening at port", zap.String("port", config.ConfigData.GrpcPort))

	if err = s.Serve(lis); err != nil {
		logger.Fatal("failed start serve", zap.Error(err))
	}
}
