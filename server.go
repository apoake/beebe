package main

import (
	"beebe/model"
	"fmt"
)

//import (
//	"beebe/controller"
//	"beebe/service"
//)
//
//func main() {
//	m := controller.Macaron()
//	m.Get("/", func() string {
//		return "Hello world!"
//	})
//	m.Run()
//	defer service.DB().Close()
//}

func main() {
	mockManager := model.GetMockManager()
	str := "@str(1, 10, lower)"
	result, err := mockManager.Mock(&str)
	fmt.Printf("%v\n%v", result, err)
}