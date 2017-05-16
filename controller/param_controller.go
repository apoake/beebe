package controller

import (
	"beebe/model"
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/session"
	"beebe/service"
	"github.com/go-macaron/binding"
)

type ParamController struct {}

func init()  {
	paramController := new(ParamController)
	Macaron().Group("/params", func() {
		Macaron().Post("save", binding.Bind(model.ParameterActionVo{}), paramController.saveParams)
	}, needLogin)
}

func (paramController *ParamController) saveParams(parameterAction model.ParameterActionVo, ctx *macaron.Context, sess session.Store) {
	if (parameterAction.ActionId == 0 && (parameterAction.RequestId != 0 || parameterAction.ResponseId != 0)) ||
		(parameterAction.ActionId > 0 && (parameterAction.ResponseId == 0 || parameterAction.RequestId == 0)) {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	userId := getCurrentUserId(sess)
	// 判断权限
	if parameterAction.ActionId > 0 {
		var projectAction *model.ProjectAction
		var ok bool
		if projectAction, ok = service.GetProjectActionService().Get(&parameterAction.ActionId); !ok {
			setErrorResponse(ctx, model.PROJECT_ACTION_NOT_FIND)
			return
		}
		if !service.GetUserService().HasProjectRight(&projectAction.ProjectId, &userId) {
			setErrorResponse(ctx, model.USER_NO_RIGHT)
			return
		}
	}
	if err := service.GetParamActionService().Save(&parameterAction); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

