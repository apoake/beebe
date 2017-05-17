package main

import (
	"beebe/controller"
	"beebe/service"
	"net/http"
	"beebe/config"
)

func main() {
	http.ListenAndServe("0.0.0.0:" + config.GetConfig().Web.Port, controller.Macaron())
	defer service.DB().Close()
}