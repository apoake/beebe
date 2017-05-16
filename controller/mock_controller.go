package controller

import (
	"gopkg.in/macaron.v1"
	"strconv"
	"beebe/service"
	"beebe/model"
)

type MockController struct {}

type MockParam struct {
	Url 		string			`json:"url"`
	ProjectId	int64 			`json:"projectId"`
}

func init() {
	mockController := new(MockController)
	Macaron().Group("/mock", func() {
		Macaron().Get("/:projectId:int/*", mockController.mock)
	})
}

func (mockController *MockController) mock(ctx *macaron.Context) {
	projectId, _ := strconv.ParseInt(ctx.Params("projectId"), 10, 64)
	url := ctx.Params("*")
	projectAction, ok := service.GetProjectActionService().GetByProjectIdAndUrl(&projectId, &url)
	if !ok {
		setErrorResponse(ctx, model.PROJECT_ACTION_NOT_FIND)
		return
	}
	result, err := service.GetMockService().MockData(&projectAction.ActionId)
	if err != nil {
		setErrorResponse(ctx, model.SYSTEM_ERROR)
		return
	}
	setSuccessResponse(ctx, result)
}
