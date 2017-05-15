package main

import (
	"beebe/service"
	"net/http"
	"beebe/controller"
)

func main() {
	http.ListenAndServe("0.0.0.0:4000", controller.Macaron())
	defer service.DB().Close()
}