package model

import (
	"testing"
	"strings"
)

func Test_convertStringArrToInfaWithA(t *testing.T) {
	str := "1, 3, @str(2, 3, @str(3, 3, \"tt\"))"
	strArr := strings.Split(str, MOCK_SPLIT)
	if !strings.Contains(str, MOCK_PREFIX) {
		arr := convertStringArrToInfa(&strArr)
		println(arr)
	}
	arr, err := convertStringArrToInfaWithA(&strArr)
	if err != nil {
		println(err)
	}
	for _, val := range *arr {
		print(val)
	}
	println(arr)
}
