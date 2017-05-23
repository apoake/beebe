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

	//str := "hello_@str(@num(2,3), @num(5,10), lower)_@str(3, 3)"
	//str := "hello_@str(1, 10, lower)_@str(3,3)"
	//str := "hello_@stre(@num(1,1), @num(3,3), <index@num(1,100)>)"
	//str := "hello_@stre(@num(2,2),@num(3,5),@str(3,@num(3,6)))"
	//result, err := model.GetMockManager().MockData(&str)
	//if err != nil {
	//	println(err)
	//}
	//fmt.Printf("%v", result)
}