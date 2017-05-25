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

const (
	DEFAULT_STR_MOCK = "@str(5, 10, lower)"
	DEFAULT_NUM_MOCK = "@num(1, 100)"
	DEFAULT_BOOL_MOCK = "@bool()"
	DEFAULT_ARRAY_MOCK = "@arr(1, 5)"
)

var DefaultMockMap map[int8]string = make(map[int8]string)

func init() {
	DefaultMockMap[model.DATA_TYPE_STRING] = DEFAULT_STR_MOCK
	DefaultMockMap[model.DATA_TYPE_BOOLEAN] = DEFAULT_BOOL_MOCK
	DefaultMockMap[model.DATA_TYPE_NUMBER] = DEFAULT_NUM_MOCK
	DefaultMockMap[model.DATA_TYPE_ARRAY_OBJECT] = DEFAULT_ARRAY_MOCK
	DefaultMockMap[model.DATA_TYPE_ARRAY_BOOLEAN] = DEFAULT_ARRAY_MOCK
	DefaultMockMap[model.DATA_TYPE_STRING] = DEFAULT_ARRAY_MOCK
	DefaultMockMap[model.DATA_TYPE_ARRAY_NUMBER] = DEFAULT_ARRAY_MOCK
}

func (mockService *MockServiceImpl) MockData(actionId *int64) (interface{}, error) {
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
	result, err := getResult(&responseArr)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getResult(parameterArr *[]model.ParameterVo) (interface{}, error) {
	if parameterArr == nil {
		return nil, nil
	}
	var err error
	resultMap := make(map[string]interface{})
	for _, parameterVo := range *parameterArr {
		if val, ok := DefaultMockMap[parameterVo.DataType]; ok && parameterVo.Expression == "" {
			parameterVo.Expression = val
		}
		switch dataType := parameterVo.DataType; {
		case dataType <= model.DATA_TYPE_BOOLEAN && dataType >= model.DATA_TYPE_STRING:
			// TODO 没有mock规则；是否设置默认规则
			if parameterVo.Expression == "" {
				resultMap[parameterVo.Identifier] = nil
			} else if resultMap[parameterVo.Identifier], err = model.GetMockManager().MockData(&parameterVo.Expression); err != nil {
				return nil, err
			}
		case dataType == model.DATA_TYPE_ARRAY_STRING:
			resultMap[parameterVo.Identifier], err = model.GetMockManager().MockDataFunc(&parameterVo.Expression, func(val interface{}) (interface{}, error) {
				if val, ok := val.(int); ok {
					return getArrayVal(val, DEFAULT_STR_MOCK)
				}
				return nil, errors.New("result must be int")
			})
			if err != nil {
				return nil, err
			}
		case dataType == model.DATA_TYPE_ARRAY_NUMBER:
			resultMap[parameterVo.Identifier], err = model.GetMockManager().MockDataFunc(&parameterVo.Expression, func(val interface{}) (interface{}, error) {
				if val, ok := val.(int); ok {
					return getArrayVal(val, DEFAULT_NUM_MOCK)
				}
				return nil, errors.New("result must be int")
			})
			if err != nil {
				return nil, err
			}
		case dataType == model.DATA_TYPE_ARRAY_BOOLEAN:
			resultMap[parameterVo.Identifier], err = model.GetMockManager().MockDataFunc(&parameterVo.Expression, func(val interface{}) (interface{}, error) {
				if val, ok := val.(int); ok {
					return getArrayVal(val, DEFAULT_BOOL_MOCK)
				}
				return nil, errors.New("result must be int")
			})
			if err != nil {
				return nil, err
			}
		case dataType == model.DATA_TYPE_ARRAY_OBJECT:
			resultMap[parameterVo.Identifier], err = model.GetMockManager().MockDataFunc(&parameterVo.Expression, func(val interface{}) (interface{}, error) {
				if val, ok := val.(int); ok {
					var erro error
					result := make([]interface{}, val, val)
					for i := 0; i < val; i++ {
						result[i], erro = getResult(parameterVo.SubParam)
						if erro != nil {
							return nil, erro
						}
					}
					return result, nil
				}
				return nil, errors.New("result must be int")
			})
			if err != nil {
				return nil, err
			}
		case dataType == model.DATA_TYPE_OBJECT:
			if resultMap[parameterVo.Identifier], err = getResult(parameterVo.SubParam); err != nil {
				return nil, err
			}
		default:
			return nil, errors.New("error mock type")
		}
	}
	return resultMap, nil
}


func getArrayVal(times int, mockStr string) (interface{}, error) {
	var erro error
	result := make([]interface{}, times, times)
	for i := 0; i < times; i++ {
		result[i], erro = model.GetMockManager().MockData(&mockStr)
		if erro != nil {
			return nil, erro
		}
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
		} else if parameterVo.Expression != "" && model.GetMockManager().IsSpecifiedAnnotation(parameterVo.Expression, model.MOCK_ARRAY) {
			resultMap[parameterVo.Identifier], err = model.GetMockManager().MockDataFunc(&parameterVo.Expression, func(val interface{}) (interface{}, error) {
				if val, ok := val.(int); ok {
					var erro error
					result := make([]interface{}, val, val)
					for i := 0; i < val; i++ {
						result[i], erro = getResultMap(parameterVo.SubParam)
						if erro != nil {
							return nil, erro
						}
					}
					return result, nil
				}
				return nil, errors.New("result must be int")
			})
			if err != nil {
				return nil, err
			}
		} else if resultMap[parameterVo.Identifier], err = getResultMap(parameterVo.SubParam); err != nil {
			return nil, err
		}
	}
	return &resultMap, nil
}



