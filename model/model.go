package model

// Base Model's definition
type Model struct {
	//CreateTime 	time.Time	`gorm:"column:create_time" json:"createTime""`
	//UpdateTime 	time.Time	`gorm:"column:update_time" json:"updateTime"`
}

type Vo struct {
	//CreateTime 	time.Time	`json:"createTime""`
	//UpdateTime 	time.Time	`json:"updateTime"`
}

type RestResult struct {
	Code 		int 		`json:"code"`
	Message		string 		`json:"message"`
	Data 		interface{}	`json:"data"`
}

func (restResult *RestResult) SetData(data interface{}) *RestResult {
	restResult.Data = data
	return restResult
}

func ConvertRestResult(errCode *ErrorCode) *RestResult {
	return &RestResult{Code: errCode.Code, Message: errCode.Msg}
}

