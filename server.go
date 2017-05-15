package main

import (
	"beebe/controller"
	"beebe/service"
	"net/http"
)

func main() {
	http.ListenAndServe("0.0.0.0:4000", controller.Macaron())
	defer service.DB().Close()
}
