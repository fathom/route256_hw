package main

import (
	"log"
	"net/http"
	"route256/libs/srvwrapper"
	"route256/loms/internal/domain"
	"route256/loms/internal/handlers/cancelorder"
	"route256/loms/internal/handlers/createorder"
	"route256/loms/internal/handlers/listorder"
	"route256/loms/internal/handlers/orderpayed"
	"route256/loms/internal/handlers/stocks"

	"github.com/julienschmidt/httprouter"
)

const port = ":8081"

//LOMS (Logistics and Order Management System)
//Сервис отвечает за учет заказов и логистику.

func main() {
	router := httprouter.New()

	businessLogic := domain.New()
	createOrderHandler := createorder.New(businessLogic)
	listOrderHandler := listorder.New()
	orderPayedHandler := orderpayed.New()
	cancelOrderHandler := cancelorder.New()
	stocksHandler := stocks.New()

	router.Handler(http.MethodPost, "/createOrder", srvwrapper.New(createOrderHandler.Handle))
	router.Handler(http.MethodPost, "/listOrder", srvwrapper.New(listOrderHandler.Handle))
	router.Handler(http.MethodPost, "/orderPayed", srvwrapper.New(orderPayedHandler.Handle))
	router.Handler(http.MethodPost, "/cancelOrder", srvwrapper.New(cancelOrderHandler.Handle))
	router.Handler(http.MethodPost, "/stocks", srvwrapper.New(stocksHandler.Handle))

	log.Println("listening http at", port)
	err := http.ListenAndServe(port, router)
	log.Fatal("cannot listen http", err)
}
