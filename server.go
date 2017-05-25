package main

import (
	"beebe/controller"
	"beebe/service"
	"beebe/config"
	"net/http"
)

func main() {
	http.ListenAndServe("0.0.0.0:" + config.GetConfig().Web.Port, controller.Macaron())
	defer service.DB().Close()
}
