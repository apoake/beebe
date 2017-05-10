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

var StrMockFeatures map[string]string
var MOCK_MAP map[string]MockType
var REGEXP *regexp.Regexp = regexp.MustCompile(`^@(.+)\((.+),(.+),(.+)\)$`)
var mockManager MockManager = MockManager{}

func init() {
	StrMockFeatures = make(map[string]string)
	StrMockFeatures[STR_FEATURE_LOWER] = 1
	StrMockFeatures[STR_FEATURE_UPPER] = 2
	StrMockFeatures[STR_FEATURE_NUMBER] = 0
	StrMockFeatures[STR_FEATURE_ALL] = 3

	MOCK_MAP = make(map[string]MockType)
	MOCK_MAP[MOCK_STRING] = &StrMock{Arr: &[][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}}
	//MOCK_MAP[MOCK_STRING_REPEAT] = MOCK_STRING_REPEAT
	//MOCK_MAP[MOCK_NUMBER] = MOCK_NUMBER
	//MOCK_MAP[MOCK_DATE] = MOCK_DATE
	//MOCK_MAP[MOCK_IMAGE] = MOCK_IMAGE
	//MOCK_MAP[MOCK_INCR] = MOCK_INCR
	//MOCK_MAP[MOCK_BOOL] = MOCK_BOOL
	//MOCK_MAP[MOCK_COLOR] = MOCK_COLOR
	//MOCK_MAP[MOCK_RGB] = MOCK_RGB
	//MOCK_MAP[MOCK_RGBA] = MOCK_RGBA
	//MOCK_MAP[MOCK_TEXT] = MOCK_TEXT
	//MOCK_MAP[MOCK_NAME] = MOCK_NAME
	//MOCK_MAP[MOCK_FIRST] = MOCK_FIRST
	//MOCK_MAP[MOCK_LAST] = MOCK_LAST
	//MOCK_MAP[MOCK_URL] = MOCK_URL
	//MOCK_MAP[MOCK_EMAIL] = MOCK_EMAIL
	//MOCK_MAP[MOCK_IP] = MOCK_IP
	//MOCK_MAP[MOCK_ADDRESS] = MOCK_ADDRESS
	//MOCK_MAP[MOCK_ZIP] = MOCK_ZIP
	//MOCK_MAP[MOCK_PCIK] = MOCK_PCIK
	//MOCK_MAP[MOCK_ARRAY] = MOCK_ARRAY


}

type MockManager struct{}

func (mockManager *MockManager) Mock(str *string) (*string, error) {
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
			val.MockVal(arr[2:])
		}
		return val.MockVal(nil)
	}
	return str, nil
}

type MockType interface {
	MockVal(params *[]string) (*interface{}, error)
	//CheckRule(str *string) error
}

type BaseMock struct {
	MockReg *regexp.Regexp
}

func (baseMock *BaseMock) Mock(str *string) (*interface{}, error) {
	if *str == "" {
		return nil, errors.New("mock str is empty")
	}
	index := strings.Index(*str, MOCK_BRACKET_LEFT)
	rs := []rune(*str)
	prefix := string(rs[:index])
	if val, ok := MOCK_MAP[prefix]; ok {
		// mock
		return val.MockVal(str)
	}
	return nil, errors.New("not support " + prefix)
}

type StrMock struct {
	Arr 	 	*[][]int
	BaseMock
}

func (strMock StrMock) MockVal(params *[]string) (*interface{}, error) {
	var err error
	var min, max int
	var features int
	if min, err = getInt(&params[0]); err != nil {
		return &params[0], err
	}
	if max, err = getInt(&params[1]); err != nil {
		return &params[1], err
	}
	if max < min {
		return nil, errors.New("max must gt min")
	}
	if length := len(params); length == 3 {
		var features string
		if features, err = getValue(&params[2]); err != nil {
			return &params[2], err
		} else if val, ok := StrMockFeatures[features]; ok {
			features = val
		} else {
			features = 3
		}
	}
	rand.Seed(time.Now().UnixNano())
	size := rand.Intn(max - min + 1) + min
	return string(strRand(size, strMock.Arr, features))
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

func (strRepeat *StrRepeat) MockVal(params *[]string) (*interface{}, error) {
	var err error
	var min, max int
	var val string
	if min, err = getInt(&params[0]); err != nil {
		return &params[0], err
	}
	if max, err = getInt(&params[1]); err != nil {
		return &params[1], err
	}
	if max < min {
		return nil, errors.New("max must gt min")
	}
	if val, err = getValue(&params[2]); err != nil {
		return nil, err
	}
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	ts := random.Intn(max - min + 1) + min
	return strings.Repeat(val, ts)
}

type NumMock struct {
	BaseMock
}

func (numMock *NumMock) MockVal(params *[]string) (*interface{}, error) {
	var err error
	var min, max, dmin, dmax int = -1, -1, -1, -1
	var intpart int = 0
	if params != nil {
		for index, val := range *params {
			var num int
			if tmp := strings.TrimSpace(val); tmp != "" {
				if num, err = getInt(&val); err != nil {
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
	if dmin == -1 {
		dmin = MIN_INT
	}
	if dmax == -1 {
		dmax = MAX_INT
	}
	intpart = random.Intn(dmax - dmin + 1) + dmin
	if dmin == -1 && dmax == -1 {
		// 返回int
		return *intpart
	} else {
		if dmin == -1 {
			dmin = 1
		}
		if dmax == -1 {
			dmax = 10
		}
		return fmt.Sprintf("%." + strconv.Itoa(dmax-dmin) + "f", float64(intpart) + random.Float64())
	}
}

type DateMock struct {
	BaseMock
}

func (dateMock *DateMock) MockVal(params *[]string) (*interface{}, error) {
	var dataForm string
	if length := len(*params); length > 0 {
		dataForm = strings.TrimSpace(*params[0])
	}
	return mockTime(dataForm, false)
}

func mockTime(dataFormat string, isMockNow bool) (*interface{}, error) {
	var err error
	dataForm := DEFAULT_DATE_FORMAT
	// 日期类型转义
	if dataFormat == "" {
		if dataForm, err = getValue(&dataFormat); err != nil {
			return nil, errors.New("error data fromat: " + dataFormat)
		}
	}
	t := time.Now()
	if !isMockNow {
		tNow := t.Unix()
		random := rand.New(rand.NewSource(t.UnixNano()))
		t = time.Unix(tNow - random.Int63n(tNow), 0)
	}
	return &(t.Format(dataForm)), nil
}

type NowMock struct {
	BaseMock
}

func (nowMock *NowMock) MockVal(params *[]string) (*interface{}, error) {
	var dataForm string
	if length := len(*params); length > 0 {
		dataForm = strings.TrimSpace(*params[0])
	}
	return mockTime(dataForm, true)
}

type ImageMock struct {
	BaseMock
}

func (imageMock *ImageMock) MockVal(params *[]string) (*interface{}, error) {
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
				if str, err = getValue(&val); err != nil {
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
	return &(fmt.Sprintf("https://dummyimage.com/%s/%s/%s/%s&text=%s", size, bcolor, fcolor, format, text))
}

type BoolMock struct {
	BaseMock
}

func (boolMock *BoolMock) MockVal(params *[]string) (*interface{}, error) {
	var truet int = 5
	if params == nil && len(params) > 0 {
		if num, err := getInt(&params[0]); err != nil {
			return nil , err
		} else if num < 10 && num > 0 {
			truet = num
		}
	}
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	if random.Intn(10) > truet {
		return false
	}
	return true
}

type ColorMock struct {
	BaseMock
}

func (colorMock *ColorMock) MockVal(params *[]string) (*interface{}, error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := "#" + fmt.Sprintf("%x", random.Intn(16777216))
	return &result, nil
}

type RgbMock struct {
	BaseMock
}

func (rgbMock *RgbMock) MockVal(params *[]string) (*interface{}, error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := fmt.Sprintf("rgb(%d, %d, %d)", random.Intn(256), random.Intn(256), random.Intn(256))
	return &result, nil
}

type RgbaMock struct {
	BaseMock
}

func (rgbaMock *RgbaMock) MockVal(params *[]string) (*interface{}, error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	result := fmt.Sprintf("rgba(%d, %d, %d, %0.2f)", random.Intn(256), random.Intn(256), random.Intn(256), float64(random.Intn(100)))
	return &result, nil
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
}

type IPMock struct {
	BaseMock
}

type AddressMock struct {
	BaseMock
}

type ZipMock struct {
	BaseMock
}

type PickMock struct {
	BaseMock
}

type IncrMock struct {
	BaseMock
}

func getValue(str *string) (string, error) {
	tmpStr := strings.TrimSpace(*str)
	if strings.Contains(tmpStr, MOCK_PREFIX) {
		return mockManager.Mock(&tmpStr)
	}
	return tmpStr, nil
}

func getInt(str *string) (int, error) {
	result, err := getValue(str)
	if err != nil {
		return nil, err
	}
	return strconv.Atoi(result), nil
}