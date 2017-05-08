package model

import (
	"strings"
	"github.com/pkg/errors"
)

const (
	MOCK_SPLIT = ","
	MOCK_PREFIX = "@"
	MOCK_BRACKET_LEFT = "("
	MOCK_BRACKET_RIGHT = ")"

	MOCK_STRING = "@str"
	MOCK_STRING_REPEAT = "@stre"
	MOCK_NUMBER = "@num"
	MOCK_DATE = "@date"
	MOCK_IMAGE = "@img"
	MOCK_INCR = "@incr"
	MOCK_BOOL = "@bool"
	MOCK_COLOR = "@color"
	MOCK_RGB = "@rgb"
	MOCK_RGBA = "@rgba"
	MOCK_TEXT = "@text"
	MOCK_NAME = "@name"
	MOCK_FIRST = "@first"
	MOCK_LAST = "@last"
	MOCK_URL = "@url"
	MOCK_EMAIL = "@email"
	MOCK_IP = "@ip"
	MOCK_ADDRESS = "@address"
	MOCK_ZIP = "@zip"
	MOCK_PCIK = "@pick"
	MOCK_ARRAY = "@arr"


	STR_FEATURE_LOWER = "lower"
	STR_FEATURE_UPPER = "upper"
	STR_FEATURE_NUMBER = "number"
	STR_FEATURE_SYMBOL = "symbol"
)

var StrMockFeatures map[string]string

func init() {
	StrMockFeatures = make(map[string]string)
	StrMockFeatures[STR_FEATURE_LOWER] = STR_FEATURE_LOWER
	StrMockFeatures[STR_FEATURE_UPPER] = STR_FEATURE_UPPER
	StrMockFeatures[STR_FEATURE_NUMBER] = STR_FEATURE_NUMBER
	StrMockFeatures[STR_FEATURE_SYMBOL] = STR_FEATURE_SYMBOL
}

type MockType interface {
	InitParams(params *[]interface{}) error
	CheckParams(params *[]interface{}) error
	Mock(params *[]interface{}) (string, error)
}

type BaseMock struct {}

func (baseMock *BaseMock) GetParams(str string) (*[]interface{}, error) {
	if str == "" {
		return nil, nil
	}
	strArr := strings.Split(str, MOCK_SPLIT)
	if !strings.Contains(str, MOCK_PREFIX) {
		return convertStringArrToInfa(&strArr), nil
	}
	return convertStringArrToInfaWithA(&strArr)
}

func (baseMock *BaseMock) GetMockType(str *string) (*MockType, error) {

}


func convertStringArrToInfa(strArr *[]string) (*[]interface{}) {
	result := make([]interface{}, 0, len(*strArr))
	for _, val := range *strArr {
		if tmp := strings.TrimSpace(val); tmp == "" {
			continue
		}
		result = append(result, strings.TrimSpace(val))
	}
	return &result
}

func convertStringArrToInfaWithA(strArr *[]string) (*[]interface{}, error) {
	println(len(*strArr))
	result := make([]interface{}, 0, len(*strArr))
	count := 0
	subStr := ""
	for _, val := range *strArr {
		var tmp string
		if tmp = strings.TrimSpace(val); tmp == "" {
			continue
		}
		if leftNum := strings.Count(tmp, MOCK_BRACKET_LEFT); leftNum > 0 {
			count += leftNum
			subStr += tmp
			continue
		} else if rightNum := strings.Count(tmp, MOCK_BRACKET_RIGHT); rightNum > 0 {
			count -= rightNum
			subStr += tmp
			if count > 0 {
				continue
			} else if count < 0 {
				return nil, errors.New("check error")
			}
			result = append(result, subStr)
			continue
		} else if count > 0 {
			subStr += tmp
			continue
		} else if count == 0 {
			result = append(result, tmp)
			continue
		}
		return nil, errors.New("check error")
	}
	return &result, nil
}


type StrMock struct {
	min 		int8
	max 		int8
	features 	string
	mockType	*MockType
	BaseMock
}

func (strMock StrMock) InitParams(params *[]interface{}) error {
	var num int
	if num = len(*params); num > 3 || num < 2 {
		return "", errors.New("params error")
	}
	strMock.min = (*params[0]).(int8)
	strMock.max = (*params[1]).(int8)
	if num == 3 {
		tmpStr := strings.TrimSpace((*params[2]).(string))
		if tmpStr == "" {
			return nil
		} else if strings.Contains(tmpStr, MOCK_PREFIX) {
			strMock.subContent = tmpStr
		} else if val, ok := StrMockFeatures[tmpStr]; ok  {
			strMock.features = val
		}
		return errors.New("third params " + tmpStr + " is invalid")
	}
	return nil
}

func (strMock StrMock) Mock(params *[]interface{}) (string, error) {

}





