package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	LomsV1 "route256/loms/internal/api/loms_v1"
	"route256/loms/internal/config"
	"route256/loms/internal/domain"
	db "route256/loms/internal/repository/db_repository"
	"route256/loms/internal/repository/db_repository/transactor"
	desc "route256/loms/pkg/loms_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// LOMS (Logistics and Order Management System)
// Сервис отвечает за учет заказов и логистику.

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", config.ConfigData.GrpcPort))
	if err != nil {
		log.Fatalf("failed start listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	dbpool, err := pgxpool.Connect(context.Background(), config.ConfigData.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	tm := transactor.NewTransactionManager(dbpool)
	ordersRepo := db.NewOrdersRepository(tm)
	orderItemsRepo := db.NewOrderItemsRepository(tm)

	businessLogic := domain.New(
		tm,
		ordersRepo,
		orderItemsRepo,
	)

	desc.RegisterLomsV1Server(s, LomsV1.NewLomsV1(businessLogic))

	log.Printf("grpc server listening at %v port", config.ConfigData.GrpcPort)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed start serve: %v", err)
	}
}
