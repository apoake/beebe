package service

import (
	"beebe/model"
	"fmt"
	"github.com/pkg/errors"
	"github.com/jinzhu/gorm"
	"encoding/json"
	"strconv"
)

var paramService = new(ParamServiceImpl)
var paramActionService = new(ParamActionServiceImpl)

func ParamService() *ParamServiceImpl {
	return paramService
}

func ParamActionService() *ParamActionServiceImpl {
	return paramActionService
}

func getDB(db *gorm.DB) *gorm.DB {
	if db == nil {
		return DB()
	}
	return db
}

//ParamServiceImpl
type ParamService interface {
	CheckParams(params *[]model.ParameterVo) error
	Delete(parameterIds *[]int64, db *gorm.DB) error
}

type ParamServiceImpl struct{}

func (paramService *ParamServiceImpl) CheckParams(params *[]model.ParameterVo) error {
	paramArr := *params
	if paramArr == nil || len(paramArr) == 0 {
		return nil
	}
	for _, param := range paramArr {
		if param.Name == "" {
			return errors.New("param[name] is empty")
		}
		if param.Identifier == "" {
			return errors.New("param[identifier] with name[" + param.Name + "] is empty")
		}
		if param.DataType < model.DATA_TYPE_STRING && param.DataType > model.DATA_TYPE_ARRAY_OBJECT {
			return errors.New("param[dataType] with name[" + param.Name + "] is empty")
		}
		if (param.DataType == model.DATA_TYPE_OBJECT || param.DataType == model.DATA_TYPE_ARRAY_OBJECT) {
			if param.SubParam == nil {
				return errors.New("param[subParam] is empty with dataType[" + fmt.Sprintf("%d", param.DataType) + "]")
			} else {
				paramService.CheckParams(param.SubParam)
			}
		}
	}
	return nil
}

func (paramService *ParamServiceImpl) Delete(parameterIds *[]int64, db *gorm.DB) error {
	return getDB(db).Where("id in (?)", parameterIds).Delete(model.ComplexParameter{}).Error
}


//ComplexParameterServiceImpl
type ComplexParameterService interface {
	DeleteByParameterIds(paramIds *[]int64, db *gorm.DB) error
	DeleteByActionId(actionId *int64, db *gorm.DB) error
	//Create(complexParameter *model.ComplexParameter) error
	//GetByParameterIds(paramIds *[]int64) (*[]model.ComplexParameter, error)
}

type ComplexParameterServiceImpl struct {}

func (complexParameterService *ComplexParameterServiceImpl) DeleteByParameterIds(paramIds *[]int64, db *gorm.DB) error {
	return getDB(db).Where("parameter_id in (?)", paramIds).Delete(model.ComplexParameter{}).Error
}

func (complexParameterService *ComplexParameterServiceImpl) DeleteByActionId(actionId *int64, db *gorm.DB) error {
	return getDB(db).Where("active_id = ?", actionId).Delete(model.ComplexParameter{}).Error
}

// ParamServiceImpl
type ParamActionServiceImpl struct{}


type ParamActionService interface {
	Get(actionId *int64) (*model.ParameterAction, error)
	Save(parameterAction *model.ParameterActionVo) error
}
// projectServiceImpl
func (paramActionService *ParamActionServiceImpl) Get(actionId *int64) (*model.ParameterAction, error) {
	paramAction := new(model.ParameterAction)
	err := DB().Find(paramAction).Error
	return paramAction, err
}

func (paramActionService *ParamActionServiceImpl) Save(parameterAction *model.ParameterActionVo) (erro error) {
	if parameterAction.ActionId < 1 {
		return errors.New("params[actionId] is empty")
	}
	projectAction, err := ProjectActionService().Get(parameterAction.ActionId)
	if err != nil {
		return err
	}
	if projectAction == nil {
		return errors.New("not find project_action by actionId: " + strconv.Itoa(parameterAction.ActionId))
	}
	requestParams := parameterAction.RequestParameter
	if err := ParamService().CheckParams(requestParams); err != nil {
		return err
	}
	responseParams := parameterAction.ResponseParameter
	if err := ParamService().CheckParams(responseParams); err != nil {
		return err
	}
	var dbParamAction *model.ParameterAction
	if parameterAction.ActionId > 0 {
		dbParamAction, erro = paramActionService.Get(parameterAction.ActionId)
	} else {
		dbParamAction = new(model.ParameterAction)
	}
	tx := DB().Begin()
	defer func() {
		if erro != nil {
			tx.Rollback()
		}
	}()
	if dbParamAction.RequestId > 0 {
		if err := tx.Where("parameter_id = ?", parameterAction.RequestId).Delete(model.ComplexParameter{}).Error; err != nil {
			return err
		}
	}
	if parameterAction.ResponseId > 0 {
		if err := tx.Where("parameter_id = ?", parameterAction.ResponseId).Delete(model.ComplexParameter{}).Error; err != nil {
			return err
		}
	}
	topParameter := &model.Parameter{Remark: model.TOP_REQUEST}
	topResponseParameter := &model.Parameter{Remark: model.TOP_RESPONSE}
	if err := tx.Save(topParameter).Error; err != nil {
		return err
	}
	if err := tx.Save(topResponseParameter).Error; err != nil {
		return err
	}
	// TODO 重构代码
	if err := saveSubParam(tx, projectAction.ActionId, topParameter.ID, requestParams); err != nil {
		return err
	}
	if err := saveSubParam(tx, projectAction.ActionId, topResponseParameter.ID, responseParams); err != nil {
		return err
	}
	requestArr, err := json.Marshal(requestParams)
	if err != nil {
		return err
	}
	responseArr, err := json.Marshal(responseParams)
	if err != nil {
		return err
	}
	parameterAction.RequestId = topParameter.ID
	parameterAction.ResponseId = topResponseParameter.ID
	parameterAction.RequestParameter = string(requestArr)
	parameterAction.ResponseParameter = string(responseArr)
	if err := tx.Save(parameterAction).Error; err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func saveSubParam(tx *gorm.DB, actionId int64, parentId int64, subParams *[]model.ParameterVo) error {
	for _, responseTmpVo := range *subParams {
		responseTmp := responseTmpVo.Convert()
		if err := tx.Save(responseTmp).Error; err != nil {
			return err
		}
		if err := tx.Save(&model.ComplexParameter{ParameterId: parentId, SubParameterId: responseTmp.ID}).Error; err != nil {
			return err
		}
		if responseTmpVo.SubParam != nil {
			if err := saveSubParam(tx, actionId, responseTmp.ID, responseTmpVo.SubParam); err != nil {
				return err
			}
		}
	}
	return nil
}




