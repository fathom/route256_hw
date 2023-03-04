package main

import (
	"fmt"
	"log"
	"net"
	"route256/loms/internal/domain"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	LomsV1 "route256/loms/internal/api/loms_v1"
	desc "route256/loms/pkg/loms_v1"
)

//LOMS (Logistics and Order Management System)
//Сервис отвечает за учет заказов и логистику.

const grpcPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", grpcPort))
	if err != nil {
		log.Fatalf("failed start listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	businessLogic := domain.New()

	desc.RegisterLomsV1Server(s, LomsV1.NewLomsV1(businessLogic))

	log.Printf("grpc server listening at %v port", grpcPort)

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed start serve: %v", err)
	}
}
