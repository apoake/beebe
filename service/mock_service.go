package service

import (
	"github.com/pkg/errors"
	"beebe/model"
	"encoding/json"
)

var mockService *MockServiceImpl = &MockServiceImpl{}

type MockService interface{
	MockData(actionId *int64) (*map[string]interface{}, error)
}

type MockServiceImpl struct {}

func GetMockService() *MockServiceImpl {
	return mockService
}


func (mockService *MockServiceImpl) MockData(actionId *int64) (*map[string]interface{}, error) {
	parameterAction, ok := GetParamActionService().Get(actionId)
	if !ok {
		return nil, errors.New("not find record")
	}
	if parameterAction.ResponseParameter == "" {
		return nil, errors.New("no response parameter")
	}
	responseArr := make([]model.ParameterVo, 5)
	if err := json.Unmarshal([]byte(parameterAction.ResponseParameter), &responseArr); err != nil {
		return nil, err
	}
	var result *map[string]interface{}
	result, err := getResultMap(&responseArr)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getResultMap(parameterArr *[]model.ParameterVo) (*map[string]interface{}, error) {
	resultMap := make(map[string]interface{})
	if parameterArr == nil {
		return nil, nil
	}
	var err error
	for _, parameterVo := range *parameterArr {
		if parameterVo.SubParam == nil {
			if parameterVo.Expression == "" {
				resultMap[parameterVo.Identifier] = nil
			} else if resultMap[parameterVo.Identifier], err = model.GetMockManager().MockData(&parameterVo.Expression); err != nil {
				return nil, err
			}
		} else if resultMap[parameterVo.Identifier], err = getResultMap(parameterVo.SubParam); err != nil {
			return nil, err
		}
	}
	return &resultMap, nil
}



