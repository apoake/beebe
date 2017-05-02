package model


type Parameter struct {
	Model
	ID             int64         `gorm:"primary_key" json:"id"`
	Name           string        `grom:"column:"name" json:"name"`
	Identifier     string        `grom:"column:"identifier" json:"identifier"`
	DataType       string        `grom:"column:"data_type" json:"dataType"`
	Remark         string        `grom:"column:"remark" json:"remark"`
	ExpressionType int8          `grom:"column:"expression_type" json:"expressionType"`
	Expression     string        `grom:"column:"expression" json:"expression"`
	MockData       string        `grom:"column:"mock_data" json:"mockData"`
}

func (Parameter) TableName() string {
	return "parameter"
}

type ParameterTemplate struct {
	Model
	ID             int64         `gorm:"primary_key" json:"id"`
	Code		   string		 `grom:"column:"code" json:"code"`
	Name           string        `grom:"column:"name" json:"name"`
	ParameterId	   int64		 `grom:"column:"parameter_id" json:"parameterId"`
}

func (ParameterTemplate) TableName() string {
	return "parameter_template"
}

type ParameterAction struct {
	Model
	ActionId 			int64		`gorm:"primary_key"`
	RequestParameter 	string		`grom:"column:request_parameter"`
	ResponseParameter 	string		`grom:"column:response_parameter"`
	RequestId 			int64		`grom:"column:request_id"`
	ResponseId  		int64 		`grom:"column:response_id"`
	ProjectId			int64		`grom:"column:project_id"`
}

func (ParameterAction) TableName() string {
	return "parameter_action"
}

type ComplexParameter struct {
	Model
	ParameterId				int64		`gorm:"column:parameter_id" json:"parameterId"`
	SubParameterId			int64		`grom:"column:sub_parameter_id" json:"subParameterId"`
}

func (ComplexParameter) TableName() string {
	return "complex_parameter"
}