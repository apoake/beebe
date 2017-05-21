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

	//str := "hello @str(1, 10, lower)"
	//result, err := model.GetMockManager().MockData(&str)
	//if err != nil {
	//	println(err)
	//}
	//fmt.Printf("%v", result)
}