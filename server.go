package main

import (
	"fmt"
	"regexp"
	"beebe/model"
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
	//result, err := mockManager.Mock(&str)
	//fmt.Printf("%v\n%v", result, err)
	println("")
	str = "@stre(1, 10, @str(1, 10, lower))"
	reg := regexp.MustCompile(`^(@.+?)\((.+)\)$`)
	arr := reg.FindStringSubmatch(str)
	for _, val :=range arr {
		println(val)
	}
	fmt.Printf("%v", arr)
	result, err := mockManager.Mock(&str)
	fmt.Printf("%v\n%v", result, err)
}