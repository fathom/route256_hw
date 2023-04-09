package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	LomsV1 "route256/loms/internal/api/loms_v1"
	"route256/loms/internal/config"
	"route256/loms/internal/domain"
	"route256/loms/internal/interceptors"
	"route256/loms/internal/kafka"
	"route256/loms/internal/logger"
	"route256/loms/internal/metrics"
	orderStauts "route256/loms/internal/notifications/order_status"
	db "route256/loms/internal/repository/db_repository"
	"route256/loms/internal/repository/db_repository/transactor"
	"route256/loms/internal/worker"
	desc "route256/loms/pkg/loms_v1"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// LOMS (Logistics and Order Management System)
// Сервис отвечает за учет заказов и логистику.

const amountWorkers = 5

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	// Инициализация логирования
	logger.Init(config.ConfigData.Dev)

	// Инициализация grpc сервера
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", config.ConfigData.GrpcPort))
	if err != nil {
		logger.Fatal("failed start listen", zap.Error(err))
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

	ctx := context.Background()

	// Инициализация базы данных
	dbpool, err := pgxpool.Connect(ctx, config.ConfigData.DatabaseURL)
	if err != nil {
		logger.Fatal("unable to create connection pool", zap.Error(err))
	}
	defer dbpool.Close()

	// Инициализация kafka продюсеров
	producer, err := kafka.NewSyncProducer(config.ConfigData.KafkaBrokers)
	if err != nil {
		logger.Fatal("unable to create kafka producer", zap.Error(err))
	}

	orderStatusSender := orderStauts.NewOrderStatusSender(producer, "orders")

	// Инициализация зависимостей, репозиториев, транзактора
	tm := transactor.NewTransactionManager(dbpool)
	ordersRepo := db.NewOrdersRepository(tm)
	orderItemsRepo := db.NewOrderItemsRepository(tm)
	warehouseRepo := db.NewWarehouseRepository(tm)

	// Инициализация воркеров снятия резервации у неоплаченных заказов
	drw, err := worker.NewDeleteReservationWorker(
		ctx,
		amountWorkers,
		ordersRepo,
		warehouseRepo,
	)
	if err != nil {
		logger.Fatal("failed start DeleteReservationWorker", zap.Error(err))
	}

	// Инициализация домена
	businessLogic := domain.New(
		tm,
		ordersRepo,
		orderItemsRepo,
		warehouseRepo,
		drw,
		orderStatusSender,
	)

	// Запуск grpc сервера
	desc.RegisterLomsV1Server(s, LomsV1.NewLomsV1(businessLogic))
	logger.Info("grpc server listening at port", zap.String("port", config.ConfigData.GrpcPort))

	grpcPrometheus.Register(s)

	go func() {
		http.Handle("/metrics", metrics.NewMetricHandler())
		err := http.ListenAndServe(":8081", nil)
		logger.Fatal("failed metrics", zap.Error(err))
		panic(err)
	}()

	if err = s.Serve(lis); err != nil {
		logger.Fatal("failed start serve", zap.Error(err))
	}
}
