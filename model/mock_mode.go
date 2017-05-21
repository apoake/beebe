package model

import (
	"strings"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"time"
	"math/rand"
	"fmt"
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
	STR_FEATURE_ALL = "all"

	MAX_INT = 1 << 32 -1
	MIN_INT = 0 - MAX_INT -1
	DEFAULT_DATE_FORMAT = "2006-01-02 15:04:05"
)

var StrMockFeatures map[string]int
var MOCK_MAP map[string]MockType
var REGEXP *regexp.Regexp = regexp.MustCompile(`^(@.+)\(((.+)[,]?)*\)$`)
var mockManager *MockManager = &MockManager{}

func init() {
	StrMockFeatures = make(map[string]int)
	StrMockFeatures[STR_FEATURE_LOWER] = 1
	StrMockFeatures[STR_FEATURE_UPPER] = 2
	StrMockFeatures[STR_FEATURE_NUMBER] = 0
	StrMockFeatures[STR_FEATURE_ALL] = 3

	MOCK_MAP = make(map[string]MockType)
	MOCK_MAP[MOCK_STRING] = &StrMock{Arr: &[][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}}
	MOCK_MAP[MOCK_STRING_REPEAT] = &StrRepeat{}
	MOCK_MAP[MOCK_NUMBER] = &NumMock{}
	MOCK_MAP[MOCK_DATE] = &DateMock{}
	MOCK_MAP[MOCK_IMAGE] = &ImageMock{}
	//MOCK_MAP[MOCK_INCR] =
	MOCK_MAP[MOCK_BOOL] = &BoolMock{}
	MOCK_MAP[MOCK_COLOR] = &ColorMock{}
	MOCK_MAP[MOCK_RGB] = &RgbMock{}
	MOCK_MAP[MOCK_RGBA] = &RgbaMock{}
	//MOCK_MAP[MOCK_TEXT] = MOCK_TEXT
	//MOCK_MAP[MOCK_NAME] = MOCK_NAME
	//MOCK_MAP[MOCK_FIRST] = MOCK_FIRST
	//MOCK_MAP[MOCK_LAST] = MOCK_LAST
	//MOCK_MAP[MOCK_URL] = MOCK_URL
	MOCK_MAP[MOCK_EMAIL] = &EmailMock{Arr: &[][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}}
	MOCK_MAP[MOCK_IP] = &IPMock{}
	//MOCK_MAP[MOCK_ADDRESS] = MOCK_ADDRESS
	MOCK_MAP[MOCK_ZIP] = &ZipMock{}
	//MOCK_MAP[MOCK_PCIK] = MOCK_PCIK
	//MOCK_MAP[MOCK_ARRAY] = MOCK_ARRAY
}

type MockManager struct{}

func GetMockManager() *MockManager {
	return mockManager
}

func (mockManager *MockManager) Mock(str *string) (interface{}, error) {
	if *str == "" {
		return nil, errors.New("mock str is empty")
	}
	arr := REGEXP.FindStringSubmatch(*str)
	length := len(arr)
	if length < 2 {
		return str, nil
	}
	mockType := arr[1]
	if val, ok := MOCK_MAP[mockType]; ok {
		if length > 2 {
			arrtt := arr[2:]
			return val.MockVal(&arrtt)
		}
		return val.MockVal(nil)
	}
	return str, nil
}

//func (mockManager *MockManager) warpResult (res interface{}, start, end string) (interface{}, error) {
//
//}

func (mockManager *MockManager) MockData(str *string) (interface{}, error) {
	if *str == "" {
		return nil, errors.New("mock str is empty")
	}
	index := strings.Index(*str, MOCK_BRACKET_LEFT)
	preIndex := strings.Index(*str, MOCK_PREFIX)
	lastIndex := strings.LastIndex(*str, MOCK_BRACKET_RIGHT)
	rs := []rune(*str)
	prefix := string(rs[preIndex:index])
	resultPrefix := string(rs[:preIndex])
	resultLast := string(rs[lastIndex+1:])
	println(resultPrefix)
	println(resultLast)
	if val, ok := MOCK_MAP[prefix]; ok {
		var mockParams *[]string
		if index + 1 == lastIndex {
			mockParams = nil
		} else {
			paramStr := string(rs[index + 1:lastIndex])
			if paramStrTrim := strings.TrimSpace(paramStr); paramStrTrim == "" {
				mockParams = nil
			} else {
				mockParams = mockManager.getMockParams(&paramStrTrim)
			}
		}
		if resultPrefix == "" && resultLast == "" {
			return val.MockVal(mockParams)
		}

	}
	return nil, errors.New("not support " + prefix)
}

func (mockManager *MockManager) getMockParams(str *string) *[]string {
	arr := strings.Split(*str, ",")
	result := make([]string, 0, len(arr))
	tmpStr := ""
	inTime := 0
	for _, val := range arr {
		str := strings.TrimSpace(val)
		if strings.Contains(str, MOCK_BRACKET_LEFT) {
			inTime++
		}
		if strings.Contains(str, MOCK_BRACKET_RIGHT) {
			inTime--
		}
		tmpStr += str
		if inTime == 0 {
			result = append(result, tmpStr)
			tmpStr = ""
		}
	}
	return &result
}

type MockType interface {
	MockVal(params *[]string) (interface{}, error)
	//CheckRule(str *string) error
}

type BaseMock struct {
	MockReg *regexp.Regexp
}

type StrMock struct {
	Arr 	 	*[][]int
	BaseMock
}

func (strMock StrMock) MockVal(params *[]string) (interface{}, error) {
	var err error
	var min, max int
	feature := 3
	pa := *params
	if min, err = getInt(pa[0]); err != nil {
		return pa, err
	}
	if max, err = getInt(pa[1]); err != nil {
		return pa, err
	}
	if max < min {
		return nil, errors.New("max must gt min")
	}
	if length := len(pa); length == 3 {
		var features string
		if features, err = getValue(pa[2]); err != nil {
			return pa[2], err
		} else if val, ok := StrMockFeatures[features]; ok {
			feature = val
		} else {
			feature = 3
		}
	}
	rand.Seed(time.Now().UnixNano())
	size := rand.Intn(max - min + 1) + min
	return string(strRand(size, feature, strMock.Arr)), nil
}

// 随机字符串
func strRand(size int, kind int, arr *[][]int) []byte {
	ikind, kinds, result := kind, *arr, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all {
			ikind = rand.Intn(len(*arr))
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

type StrRepeat struct {
	BaseMock
}

func (strRepeat StrRepeat) MockVal(params *[]string) (interface{}, error) {
	var err error
	var min, max int
	var val string
	pa := *params
	if min, err = getInt(pa[0]); err != nil {
		return pa[0], err
	}
	if max, err = getInt(pa[1]); err != nil {
		return pa[1], err
	}
	if max < min {
		return nil, errors.New("max must gt min")
	}
	if val, err = getValue(pa[2]); err != nil {
		return nil, err
	}
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	ts := random.Intn(max - min + 1) + min
	return strings.Repeat(val, ts), nil
}

type NumMock struct {
	BaseMock
}

func (numMock NumMock) MockVal(params *[]string) (interface{}, error) {
	var err error
	var min, max, dmin, dmax int = -1, -1, -1, -1
	var intpart int = 0
	if params != nil {
		for index, val := range *params {
			var num int
			if tmp := strings.TrimSpace(val); tmp != "" {
				if num, err = getInt(val); err != nil {
					return nil, err
				}
			}
			switch index {
			case 0:
				min = num
			case 1:
				max = num
			case 2:
				dmin = num
			case 3:
				dmax = num
			}
		}
	}
	if max < min || dmax < dmin {
		return nil, errors.New("params error")
	}
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	intpart = random.Intn(max - min + 1) + dmin
	if dmin == -1 && dmax == -1 {
		// 返回int
		return intpart, nil
	} else {
		if dmin == -1 {
			dmin = 1
		}
		if dmax == -1 {
			dmax = 10
		}
		return fmt.Sprintf("%." + strconv.Itoa(dmax-dmin) + "f", float64(intpart) + random.Float64()), nil
	}
}

type DateMock struct {
	BaseMock
}

func (dateMock DateMock) MockVal(params *[]string) (interface{}, error) {
	var dataForm string
	arr := *params
	if length := len(arr); length > 0 {
		str, _ := getValue(arr[0])
		dataForm = strings.TrimSpace(str)
	}
	return mockTime(dataForm, false)
}

func mockTime(dataFormat string, isMockNow bool) (interface{}, error) {
	var err error
	dataForm := DEFAULT_DATE_FORMAT
	// 日期类型转义
	if dataFormat == "" {
		if dataForm, err = getValue(dataFormat); err != nil {
			return nil, errors.New("error data fromat: " + dataFormat)
		}
	}
	t := time.Now()
	if !isMockNow {
		tNow := t.Unix()
		random := rand.New(rand.NewSource(t.UnixNano()))
		t = time.Unix(tNow - random.Int63n(tNow), 0)
	}
	return t.Format(dataForm), nil
}

type NowMock struct {
	BaseMock
}

func (nowMock NowMock) MockVal(params *[]string) (interface{}, error) {
	var dataForm string
	arr := *params
	if length := len(arr); length > 0 {
		str, _ := getValue(arr[0])
		dataForm = strings.TrimSpace(str)
	}
	return mockTime(dataForm, true)
}

type ImageMock struct {
	BaseMock
}

func (imageMock ImageMock) MockVal(params *[]string) (interface{}, error) {
	size := "600x400"
	bcolor := "000"
	fcolor := "fff"
	format := "png"
	text := "text"
	var err error
	if params != nil {
		for index, val := range *params {
			var str string
			if tmp := strings.TrimSpace(val); tmp != "" {
				if str, err = getValue(val); err != nil {
					return nil, err
				}
			}
			switch index {
			case 0:
				size = str
			case 1:
				bcolor = str
			case 2:
				fcolor = str
			case 3:
				format = str
			case 4:
				text = str
			}
		}
	}
	return fmt.Sprintf("https://dummyimage.com/%s/%s/%s/%s&text=%s", size, bcolor, fcolor, format, text), nil
}

type BoolMock struct {
	BaseMock
}

func (boolMock BoolMock) MockVal(params *[]string) (interface{}, error) {
	var truet int = 5
	arr := *params
	if params == nil && len(arr) > 0 {
		if num, err := getInt(arr[0]); err != nil {
			return nil , err
		} else if num < 10 && num > 0 {
			truet = num
		}
	}
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	if random.Intn(10) > truet {
		return false, nil
	}
	return true, nil
}

type ColorMock struct {
	BaseMock
}

func (colorMock ColorMock) MockVal(params *[]string) (interface{}, error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := "#" + fmt.Sprintf("%x", random.Intn(16777216))
	return result, nil
}

type RgbMock struct {
	BaseMock
}

func (rgbMock RgbMock) MockVal(params *[]string) (interface{}, error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := fmt.Sprintf("rgb(%d, %d, %d)", random.Intn(256), random.Intn(256), random.Intn(256))
	return result, nil
}

type RgbaMock struct {
	BaseMock
}

func (rgbaMock RgbaMock) MockVal(params *[]string) (interface{}, error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := fmt.Sprintf("rgba(%d, %d, %d, %0.2f)", random.Intn(256), random.Intn(256), random.Intn(256), float64(random.Intn(100)))
	return result, nil
}

type TextMock struct {
	BaseMock
}

type NameMock struct {
	BaseMock
}

type FirstMock struct {
	BaseMock
}

type LastMock struct {
	BaseMock
}

type UrlMock struct {
	BaseMock
}

type EmailMock struct {
	BaseMock
	Arr 	 	*[][]int
}

func (emailMock EmailMock) MockVal(params *[]string) (interface{}, error) {
	return fmt.Sprintf("%s.%s@%s.%s", string(strRand(2, 1, emailMock.Arr)),
		strRand(10, 1, emailMock.Arr), strRand(6, 1, emailMock.Arr),
		strRand(5, 1, emailMock.Arr)), nil
}

type IPMock struct {
	BaseMock
}

func (ipMock IPMock) MockVal(params *[]string) (interface{}, error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := fmt.Sprintf("%d.%d.%d.%d", random.Intn(256), random.Intn(256), random.Intn(256), random.Intn(256))
	return result, nil
}

type AddressMock struct {
	BaseMock
}

type ZipMock struct {
	BaseMock
}

func (zipMock ZipMock) MockVal(params *[]string) (interface{}, error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return random.Intn(899099) + 100000, nil
}

type PickMock struct {
	BaseMock
}

type IncrMock struct {
	BaseMock
}


func getValue(str string) (string, error) {
	tmpStr := strings.TrimSpace(str)
	if strings.Contains(tmpStr, MOCK_PREFIX) {
		intr, err := mockManager.Mock(&tmpStr)
		if err != nil {
			return "", err
		}
		if val, ok := intr.(string); ok {
			return val, nil
		}
		return "", errors.New("not string type")
	}
	return tmpStr, nil
}

func getInt(str string) (int, error) {
	result, err := getValue(str)
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(result)
}