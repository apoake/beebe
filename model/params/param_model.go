package params

import "beebe/model"

type Parameter struct {
	model.Model
	ID             int64         `gorm:"primary_key" json:"id"`
	Name           string        `grom:"column:"name" json:"name"`
	Identifier     string        `grom:"column:"identifier" json:"identifier"`
	DataType       string        `grom:"column:"data_type" json:"dataType"`
	Remark         string        `grom:"column:"remark" json:"remark"`
	ExpressionType int8          `grom:"column:"expression_type" json:"expressionType"`
	Expression     string        `grom:"column:"expression" json:"expression"`
	MockData       string        `grom:"column:"mock_data" json:"mockData"`
}

type ParameterTemplate struct {
	model.Model
	ID             int64         `gorm:"primary_key" json:"id"`
	Code		   string		 `grom:"column:"code" json:"code"`
	Name           string        `grom:"column:"name" json:"name"`
}