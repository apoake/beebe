package main

import (
	"beebe/service"
	"beebe/controller"
)

func main() {
	m := controller.Macaron()
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Run()
	defer service.DB().Close()

	//data, _ := ioutil.ReadFile("config/mock-min.js")
	//runtime := otto.New()
	//if _, err := runtime.Run(data); err != nil {
	//	panic(err)
	//}
	//
	//value, erro :=runtime.Run(`
	//	Mock.mock({
	//	"string|1-10": "★"
	//	})
	//`)
	//
	//if erro != nil {
	//	println(erro)
	//}
	//
	//value.Object().
	//
	//result, _ := value.Object().Get("string")
	//
	//tt, _ := result.ToString()
	//println(tt)

}
