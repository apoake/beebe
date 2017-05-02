package service

import (
	"beebe/model"
)

type ParamService interface {

}

type ParamActionService interface {
}

type ParamServiceImpl struct{}

type ParamActionServiceImpl struct{}

// ParamServiceImpl


// projectServiceImpl
func (paramActionService *ParamActionServiceImpl) GetAllByProjectId(projectId *int64,
	start *int64, limit *int64) (*[]model.ParameterAction, error) {
	parameterActions := make([]model.ParameterAction, *limit)
	err := DB().Offset(start).Limit(limit).Where("project_id = ?", projectId).Find(parameterActions).Error
	return parameterActions, err
}


