package main

import (
	"log"
	"net/http"
)

const port = ":8082"

//Notifications
//Будет слушать Кафку и отправлять уведомления, внешнего API нет.

func main() {
	log.Println("listening http at", port)
	err := http.ListenAndServe(port, nil)
	log.Fatal("cannot listen http", err)
}
