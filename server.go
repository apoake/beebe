package main

import (
)
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
}

