package main

import (
	"log"
	"net/http"

	"./handle"
)

func main() {
	http.HandleFunc("/order", handle.Order)
	http.HandleFunc("/book", handle.Book)
	err := http.ListenAndServeTLS("127.0.0.1:8081", "./cert/ca.crt", "./cert/ca.key", nil)
	log.Println("err:", err)
	return
}
