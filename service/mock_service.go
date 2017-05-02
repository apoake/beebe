package service

var MockService interface{
	MockData(actionId *int64) *map[string]interface{}
}


