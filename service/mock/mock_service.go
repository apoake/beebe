package mock

var MockService interface{
	MockData(actionId *int64) *map[string]interface{}
}


