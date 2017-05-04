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
	//var start int64 = 0
	//var limit int64 = 10
	//project, err :=service.GetProjectService().GetProjectsPage("pro", &start, &limit)
	//if err != nil {
	//	println(err)
	//}
	//println(project)
}

