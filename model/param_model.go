package model

const (
	DATA_TYPE_STRING = iota
	DATA_TYPE_NUMBER
	DATA_TYPE_OBJECT
	DATA_TYPE_BOOLEAN
	DATA_TYPE_ARRAY_STRING
	DATA_TYPE_ARRAY_NUMBER
	DATA_TYPE_ARRAY_BOOLEAN
	DATA_TYPE_ARRAY_OBJECT

	TOP_REQUEST = "request"
	TOP_RESPONSE = "response"
)

type Parameter struct {
	Model
	ID             int64         `gorm:"primary_key" json:"id"`
	Name           string        `grom:"column:"name" json:"name"`
	Identifier     string        `grom:"column:"identifier" json:"identifier"`
	DataType       int8        `grom:"column:"data_type" json:"dataType"`
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
}

func (ParameterAction) TableName() string {
	return "parameter_action"
}

type ComplexParameter struct {
	Model
	ParameterId				int64		`gorm:"column:parameter_id" json:"parameterId"`
	SubParameterId			int64		`grom:"column:sub_parameter_id" json:"subParameterId"`
	ActionId				int64		`grom:"column:action_id" json:"actionId"`
}

func (ComplexParameter) TableName() string {
	return "complex_parameter"
}

// vo
type ParameterVo struct {
	Vo
	ID			int64			`json:"id"`
	Name 		string			`json:"name"`
	Identifier 	string			`json:"identifier"`
	DataType	int8			`json:"dataType"`
	Remark		string			`json:"remark"`
	Expression	string			`json:"expression"`
	MockData	string			`json:"mockData"`
	SubParam	*[]ParameterVo	`json:"subParam"`
}

func (parameterVo *ParameterVo) Convert() *Parameter {
	return &Parameter{
		ID: parameterVo.ID,
		Name: parameterVo.Name,
		Identifier: parameterVo.Identifier,
		DataType: parameterVo.DataType,
		Remark: parameterVo.Remark,
		Expression: parameterVo.Expression,
		MockData: parameterVo.MockData}
}

type ParameterActionVo struct {
	Vo
	ActionId 			int64			`json:"actionId"`
	RequestParameter 	*[]ParameterVo	`json:"requestParameter"`
	ResponseParameter 	*[]ParameterVo	`json:"responseParameter"`
	RequestId 			int64			`json:"requestId"`
	ResponseId  		int64 			`json:"responseId"`
}