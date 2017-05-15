package main

import (
	"beebe/service"
	"net/http"
	"beebe/controller"
	"beebe/config"
)

func main() {
	http.ListenAndServe("0.0.0.0:" + config.GetConfig().Web.Port, controller.Macaron())
	defer service.DB().Close()
}