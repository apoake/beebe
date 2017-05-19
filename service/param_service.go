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
var complexParameterService = new(ComplexParameterServiceImpl)

func GetParamService() *ParamServiceImpl {
	return paramService
}

func GetParamActionService() *ParamActionServiceImpl {
	return paramActionService
}

func GetComplexParameterService() *ComplexParameterServiceImpl {
	return complexParameterService
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
	GetByActionId(actionId *int64, db *gorm.DB) (*[] model.ComplexParameter, error)
}

type ComplexParameterServiceImpl struct {}

func (complexParameterService *ComplexParameterServiceImpl) DeleteByParameterIds(paramIds *[]int64, db *gorm.DB) error {
	return getDB(db).Where("parameter_id in (?)", paramIds).Delete(model.ComplexParameter{}).Error
}

func (complexParameterService *ComplexParameterServiceImpl) DeleteByActionId(actionId *int64, db *gorm.DB) error {
	return getDB(db).Where("active_id = ?", actionId).Delete(model.ComplexParameter{}).Error
}

func (complexParameterService *ComplexParameterServiceImpl) GetByActionId(actionId *int64, db *gorm.DB) (*[] model.ComplexParameter, error) {
	complexArr := make([]model.ComplexParameter, 0)
	err := getDB(db).Where("action_id = ?", actionId).Find(&complexArr).Error
	return &complexArr, err
}

// ParamServiceImpl
type ParamActionServiceImpl struct{}


type ParamActionService interface {
	Get(actionId *int64) (*model.ParameterAction, bool)
	Save(parameterAction *model.ParameterActionVo) error
	DeleteByActionId(actionId *int64, db *gorm.DB) error
}
// projectServiceImpl
func (paramActionService *ParamActionServiceImpl) Get(actionId *int64) (*model.ParameterAction, bool) {
	paramAction := &model.ParameterAction{ActionId: *actionId}
	isExist := !DB().Find(paramAction).RecordNotFound()
	return paramAction, isExist
}

func (paramActionService *ParamActionServiceImpl) Save(parameterAction *model.ParameterActionVo) (erro error) {
	if parameterAction.ActionId < 1 {
		return errors.New("params[actionId] is empty")
	}
	projectAction, ok := GetProjectActionService().Get(&parameterAction.ActionId)
	if  !ok {
		return errors.New("not find projectAction")
	}
	// TODO: 权限限制
	if projectAction == nil {
		return errors.New("not find project_action by actionId: " + strconv.FormatInt(parameterAction.ActionId, 10))
	}
	requestParams := parameterAction.RequestParameter
	if err := GetParamService().CheckParams(requestParams); err != nil {
		return err
	}
	responseParams := parameterAction.ResponseParameter
	if err := GetParamService().CheckParams(responseParams); err != nil {
		return err
	}
	var dbParamAction *model.ParameterAction
	if parameterAction.ActionId > 0 {
		dbParamAction, ok = paramActionService.Get(&parameterAction.ActionId)
		if !ok {
			return errors.New("not find record")
		}
	} else {
		dbParamAction = new(model.ParameterAction)
	}
	tx := DB().Begin()
	defer func() {
		if erro != nil {
			tx.Rollback()
		}
	}()
	// cleam complex_parameter
	if err := GetComplexParameterService().DeleteByActionId(&parameterAction.ActionId, tx); err != nil {
		return err
	}
	topParameter := &model.Parameter{Remark: model.TOP_REQUEST}
	topResponseParameter := &model.Parameter{Remark: model.TOP_RESPONSE}
	topParams := []model.Parameter{*topParameter, *topResponseParameter}
	if err := tx.Save(&topParams).Error; err != nil {
		return err
	}
	//if err := tx.Save(topParameter).Error; err != nil {
	//	return err
	//}
	//if err := tx.Save(topResponseParameter).Error; err != nil {
	//	return err
	//}
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
	dbParamAction.RequestId = topParameter.ID
	dbParamAction.ResponseId = topResponseParameter.ID
	dbParamAction.RequestParameter = string(requestArr)
	dbParamAction.ResponseParameter = string(responseArr)
	if err := tx.Save(dbParamAction).Error; err != nil {
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

func (paramActionService *ParamActionServiceImpl) DeleteByActionId(actionId *int64, db *gorm.DB) error {
	database := getDB(db)
	comArr, err := GetComplexParameterService().GetByActionId(actionId, database)
	if err != nil {
		return err
	}
	parameterIds := getParameterIds(comArr)
	if len(*parameterIds) > 0 {
		if err := GetParamService().Delete(parameterIds, database); err != nil {
			return err
		}
	}
	if err := GetComplexParameterService().DeleteByActionId(actionId, database); err != nil {
		return err
	}
	if err := database.Where("action_id = ?", actionId).Delete(&model.ParameterAction{}).Error; err != nil {
		return err
	}
	return nil
}


func getParameterIds(complexParams *[]model.ComplexParameter) *[]int64 {
	arr := *complexParams
	comMap := make(map[int64] int64)
	result := make([]int64, 0, len(arr) + 2)
	for _, complexParam := range arr {
		if _, ok := comMap[complexParam.ParameterId]; !ok {
			comMap[complexParam.ParameterId] = complexParam.ParameterId
			result = append(result, complexParam.ParameterId)
		}
		if _, ok:= comMap[complexParam.SubParameterId]; !ok {
			comMap[complexParam.SubParameterId] = complexParam.SubParameterId
			result = append(result, complexParam.SubParameterId)
		}
	}
	return &result
}

