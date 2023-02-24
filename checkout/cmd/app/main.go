package main

import (
	"log"
	"net/http"
	"route256/checkout/internal/clients/loms"
	"route256/checkout/internal/clients/productservice"
	"route256/checkout/internal/config"
	"route256/checkout/internal/domain"
	"route256/checkout/internal/handlers/addtocart"
	"route256/checkout/internal/handlers/deletefromcart"
	"route256/checkout/internal/handlers/listcart"
	"route256/checkout/internal/handlers/purchase"
	"route256/libs/srvwrapper"

	"github.com/julienschmidt/httprouter"
)

const port = ":8080"

//Checkout
//Сервис отвечает за корзину и оформление заказа.

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal("config init", err)
	}

	router := httprouter.New()

	lomsClient := loms.New(config.ConfigData.Services.Loms)

	productServiceClient := productservice.New(config.ConfigData.Services.Product, config.ConfigData.Token)
	businessLogic := domain.New(lomsClient, productServiceClient)

	addToCartHandler := addtocart.New(businessLogic)
	deleteFromCartHandler := deletefromcart.New()
	listCartHandler := listcart.New(businessLogic)
	purchaseHandler := purchase.New(businessLogic)

	router.Handler(http.MethodPost, "/addToCart", srvwrapper.New(addToCartHandler.Handle))
	router.Handler(http.MethodPost, "/deleteFromCart", srvwrapper.New(deleteFromCartHandler.Handle))
	router.Handler(http.MethodPost, "/listCart", srvwrapper.New(listCartHandler.Handle))
	router.Handler(http.MethodPost, "/purchase", srvwrapper.New(purchaseHandler.Handle))

	log.Println("listening http at", port)
	err = http.ListenAndServe(port, router)
	log.Fatal("cannot listen http", err)
}
