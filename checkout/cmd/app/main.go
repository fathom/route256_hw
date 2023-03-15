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
	db "route256/checkout/internal/repository/db_repository"
	"route256/checkout/internal/repository/db_repository/transactor"
	desc "route256/checkout/pkg/checkout_v1"

	"github.com/jackc/pgx/v4/pgxpool"
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
	cartRepo := db.NewCartRepository(tm)

	lomsConn, err := grpc.Dial(config.ConfigData.Services.Loms, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed connect to server: %v", err)
	}
	defer lomsConn.Close()

	lomsClient := LomsClient.New(lomsConn)

	productConn, err := grpc.Dial(config.ConfigData.Services.Product, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed connect to server: %v", err)
	}
	defer productConn.Close()

	productServiceClient := ProductClient.New(productConn, config.ConfigData.Token)

	businessLogic := domain.New(
		lomsClient,
		productServiceClient,
		tm,
		cartRepo,
	)

	desc.RegisterCheckoutV1Server(s, CheckoutV1.NewCheckoutV1(businessLogic))
	log.Printf("grpc server listening at %v port", config.ConfigData.GrpcPort)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed start serve: %v", err)
	}
}
