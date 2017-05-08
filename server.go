package main

import (
	"beebe/controller"
	"beebe/service"
)

func main() {
	m := controller.Macaron()
	m.Get("/", func() string {
		return "Hello world!"
	})
	m.Run()
	defer service.DB().Close()
}

