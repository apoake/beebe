package main

import (
)
import (
	"beebe/service"
	"beebe/model"
)

func main() {
	//m := controller.Macaron()
	//m.Get("/", func() string {
	//	return "Hello world!"
	//})
	//m.Run()
	//defer service.DB().Close()


	test()

}

func test() {
	arr := make([]model.Parameter, 2)
	arr[0] = model.Parameter{DataType:"1", Expression:"1", ExpressionType:"1", Identifier:"1", Name:"1"}
	arr[1] = model.Parameter{DataType:"2", Expression:"2", ExpressionType:"2", Identifier:"2", Name:"2"}
	if err := service.DB().Create(&arr).Error; err != nil  {
		println(err)
	}
}
