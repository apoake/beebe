package mock

import (
	"time"
	"fmt"
	"strings"
	"strconv"
	"math/rand"
	"encoding/json"
)

const (
	NUMBER_MIN_ONLY = iota
	NUMBER_MIN_MAX
	NUMBER_D_NO
	NUMBER_D_MIN_ONLY
	NUMBER_D_MIN_MAX
	String string = "string"
	Number string = "number"
	Boolean string = "boolean"
	Object string = "object"
	Array string = "array"
)

type MockManager struct {}

func (mockManager *MockManager) Mock(mockType *string, val *string, value *string) interface{} {
	var mock Mock
	tepl := &Template{mockType: *mockType, val: *val, value: *value}
	switch *mockType {
	case String :
		if strings.Contains(*val, "-") {
			mock = &StringRangeTemplate{RangeTemplate: newRangeTemplate(tepl)}
		} else {
			mock = &StringCountTemplate{CountTemplate: newCountTemplate(tepl)}
		}
	case Number :
		if strings.Contains(*val, "-") {
			mock = &NumberRangeTemplate{RangeTemplate: &RangeTemplate{Template: tepl}}
		}
	case Boolean :
		if strings.Contains(*val, "-") {
			mock = &BooleanRangeTemplate{RangeTemplate:  newRangeTemplate(tepl)}
		} else {
			mock = &BooleanCountTemplate{CountTemplate: newCountTemplate(tepl)}
		}
	case Object :
		if strings.Contains(*val, "-") {
			mock = &ObjectRangeTemplate{RangeTemplate:  newRangeTemplate(tepl)}
		} else {
			mock = &ObjectCountTemplate{CountTemplate: newCountTemplate(tepl)}
		}
	case Array :
		if strings.Contains(*val, "-") {
			mock = &ArrayRangeTemplate{RangeTemplate:  newRangeTemplate(tepl)}
		} else {
			mock = &ArrayCountTemplate{CountTemplate: newCountTemplate(tepl)}
		}
	}
	if mock == nil {
		return  nil
	}
	//mock.Init()
	return mock.Mock()
}

type Mock interface {
	Mock() interface{}
	//Init()
}

type Template struct {
	mockType	string
	val			string
	value 		string
	isInit		bool
}

type RangeTemplate struct {
	*Template
	min 		int
	max 		int
}

type CountTemplate struct {
	*Template
	count		int
}

func newRangeTemplate(template *Template) *RangeTemplate {
	rangeTemplate := &RangeTemplate{Template: template}
	rangeTemplate.Init()
	return rangeTemplate
}

func (rangeTemplate *RangeTemplate) Init() {
	if rangeTemplate.isInit {
		return
	}
	rangeTemplate.isInit = true
	rangeTemplate.initValue()
}

func (rangeTemplate *RangeTemplate) initValue() {
	strArr := strings.Split(rangeTemplate.val, "-")
	min, _ := strconv.ParseInt(strArr[0], 10, 8)
	max, _ := strconv.ParseInt(strArr[1], 10, 8)
	rangeTemplate.min = int(min)
	rangeTemplate.max = int(max)
}

func newCountTemplate(template *Template) *CountTemplate {
	countTemplate := &CountTemplate{Template: template}
	countTemplate.Init()
	return countTemplate
}

func (countTemplate *CountTemplate) Init() {
	if countTemplate.isInit {
		return
	}
	countTemplate.isInit = true
	count, _ := strconv.ParseInt(countTemplate.val, 10, 8)
	countTemplate.count = int(count)
}

// 字符串范围模板
type StringRangeTemplate struct {
	*RangeTemplate
}

func (stringRangeTemplate *StringRangeTemplate) Init() {
	stringRangeTemplate.RangeTemplate.Init()
}

func (stringRangeTemplate *StringRangeTemplate) Mock() interface{} {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	ts := random.Intn(stringRangeTemplate.max - stringRangeTemplate.min + 1) + stringRangeTemplate.min
	return strings.Repeat(stringRangeTemplate.value, ts)
}

// 字符串数量模板
type StringCountTemplate struct {
	*CountTemplate
}

func (stringCountTemplate *StringCountTemplate) Init() {
	stringCountTemplate.CountTemplate.Init()
}

func (stringCountTemplate *StringCountTemplate) Mock() interface{} {
	return strings.Repeat(stringCountTemplate.value, stringCountTemplate.count)
}

// 数字范围模板
type NumberRangeTemplate struct {
	*RangeTemplate
	dMin		int
	dMax		int
	intType 	int
	decimalType int
}

// 初始化
func (numberRangeTemplate *NumberRangeTemplate) Init() {
	if numberRangeTemplate.isInit {
		return
	}
	numberRangeTemplate.isInit = true
	strArr := strings.Split(numberRangeTemplate.val, ".")
	interSArr := strings.Split(strArr[0], "-")
	if len(interSArr) == 1 {
		min, _ := strconv.ParseInt(interSArr[0], 10, 8)
		numberRangeTemplate.min = int(min)
		numberRangeTemplate.intType = NUMBER_MIN_ONLY
	} else {
		min, _ := strconv.ParseInt(interSArr[0], 10, 8)
		numberRangeTemplate.min = int(min)
		max, _ := strconv.ParseInt(interSArr[1], 10, 8)
		numberRangeTemplate.max = int(max)
		numberRangeTemplate.intType = NUMBER_MIN_MAX
	}
	numberRangeTemplate.decimalType = NUMBER_D_NO
	if len(strArr) == 2 {
		decimalArr := strings.Split(strArr[1], "-")
		if len(decimalArr) == 1 {
			dMin, _ := strconv.ParseInt(decimalArr[0], 10, 8)
			numberRangeTemplate.dMin = int(dMin)
			numberRangeTemplate.decimalType = NUMBER_D_MIN_ONLY
		} else {
			dMin, _ := strconv.ParseInt(decimalArr[0], 10, 8)
			numberRangeTemplate.dMin = int(dMin)
			dMax, _ := strconv.ParseInt(decimalArr[1], 10, 8)
			numberRangeTemplate.dMax = int(dMax)
			numberRangeTemplate.decimalType = NUMBER_D_MIN_MAX
		}
	}
}

func (numberRangeTemplate * NumberRangeTemplate) Mock() interface{} {
	var intPart int = 0
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	if numberRangeTemplate.intType == NUMBER_MIN_ONLY {
		intPart = numberRangeTemplate.min
	} else if numberRangeTemplate.intType == NUMBER_MIN_MAX {
		intPart = random.Intn(numberRangeTemplate.max - numberRangeTemplate.min + 1) + numberRangeTemplate.min
	}
	if numberRangeTemplate.decimalType == NUMBER_D_NO {
		return intPart
	} else if numberRangeTemplate.decimalType == NUMBER_D_MIN_ONLY {
		return fmt.Sprintf("%." + strconv.Itoa(numberRangeTemplate.dMin) + "f", float64(intPart) + random.Float64())
	} else if numberRangeTemplate.decimalType == NUMBER_D_MIN_MAX {
		length := random.Intn(numberRangeTemplate.dMax - numberRangeTemplate.dMin + 1) + numberRangeTemplate.dMin
		return fmt.Sprintf("%." + strconv.Itoa(length) + "f", float64(intPart) + random.Float64())
	}
	return intPart
}

type BooleanCountTemplate struct {
	*CountTemplate
}

func (booleanCountTemplate *BooleanCountTemplate) Init() {
	booleanCountTemplate.CountTemplate.Init()
}

func (booleanCountTemplate *BooleanCountTemplate) Mock() interface{} {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	randNum := random.Intn(10);
	return randNum < 5
}

type BooleanRangeTemplate struct {
	*RangeTemplate
}

func (booleanRangeTemplate *BooleanRangeTemplate) Init() {
	booleanRangeTemplate.RangeTemplate.Init()
}

func (booleanRangeTemplate *BooleanRangeTemplate) Mock() interface{} {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	randNum := random.Intn(10);
	return randNum >= booleanRangeTemplate.min && randNum <= booleanRangeTemplate.max
}

// 对象数量模板
type ObjectCountTemplate struct {
	*CountTemplate
}

func (objectCountTemplate *ObjectCountTemplate) Init() {
	objectCountTemplate.CountTemplate.Init()
}

func (objectCountTemplate *ObjectCountTemplate) Mock() interface{} {
	return getObjectResult(&objectCountTemplate.value, objectCountTemplate.count)
}

// 对象范围模板
type ObjectRangeTemplate struct {
	*RangeTemplate
}

func (objectRangeTemplate *ObjectRangeTemplate) Init() {
	objectRangeTemplate.RangeTemplate.Init()
}

func (objectRangeTemplate *ObjectRangeTemplate) Mock() interface{} {
	return getObjectResult(&objectRangeTemplate.value,
		rand.New(rand.NewSource(time.Now().UnixNano())).Intn(objectRangeTemplate.max - objectRangeTemplate.min + 1) + objectRangeTemplate.min)
}

func getObjectResult(str *string, ts int) *map[string]interface{} {
	tmpMap := make(map[string] interface{})
	if err := json.Unmarshal([]byte(*str), &tmpMap); err != nil {
		// logger error
		fmt.Println(err)
		return nil
	}
	dataLength := len(tmpMap)
	if ts >= dataLength {
		return &tmpMap
	}
	resultMap := make(map[string] interface{}, dataLength)
	index := 0
	for key, value := range tmpMap {
		if index >= dataLength {
			break
		}
		resultMap[key] = value
	}
	return &resultMap
}


type ArrayCountTemplate struct {
	*CountTemplate
}

func (arrayCountTemplate *ArrayCountTemplate) Init() {
	arrayCountTemplate.CountTemplate.Init()
}

func (arrayCountTemplate *ArrayCountTemplate) Mock() interface{} {
	return getArrayResult(&arrayCountTemplate.value, arrayCountTemplate.count)
}

type ArrayRangeTemplate struct {
	*RangeTemplate
}

func (arrayRangeTemplate *ArrayRangeTemplate) Init() {
	arrayRangeTemplate.RangeTemplate.Init()
}

func (arrayRangeTemplate *ArrayRangeTemplate) Mock() interface{} {
	return getArrayResult(&arrayRangeTemplate.value,
		rand.New(rand.NewSource(time.Now().UnixNano())).Intn(arrayRangeTemplate.max - arrayRangeTemplate.min + 1) + arrayRangeTemplate.min)
}

func getArrayResult(str *string, count int) *[]interface{} {
	tmpList := make([]interface{}, 5)
	if err := json.Unmarshal([]byte(*str), &tmpList); err != nil {
		// logger error
		fmt.Println(err)
		return nil
	}
	tmpListLength := len(tmpList)
	dataLength := tmpListLength * count
	resultList := make([]interface{}, dataLength)
	tIndex := 0
	for index:= 0; index < dataLength; index++ {
		if tIndex >= tmpListLength {
			tIndex = tIndex % tmpListLength
		}
		resultList[index] = tmpList[tIndex]
		tIndex++
	}
	return &resultList
}
