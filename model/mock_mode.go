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
	MOCK_SPLIT         = ","
	MOCK_PREFIX        = "@"
	MOCK_BRACKET_LEFT  = "("
	MOCK_BRACKET_RIGHT = ")"

	MOCK_STRING        = "@str"
	MOCK_STRING_REPEAT = "@stre"
	MOCK_NUMBER        = "@num"
	MOCK_DATE          = "@date"
	MOCK_IMAGE         = "@img"
	//MOCK_INCR          = "@incr"
	MOCK_BOOL          = "@bool"
	MOCK_COLOR         = "@color"
	MOCK_RGB           = "@rgb"
	MOCK_RGBA          = "@rgba"
	MOCK_NAME          = "@name"
	MOCK_FIRST         = "@first"
	MOCK_LAST          = "@last"
	MOCK_URL           = "@url"
	MOCK_EMAIL         = "@email"
	MOCK_IP            = "@ip"
	MOCK_ADDRESS       = "@address"
	MOCK_REGION		   = "@region"
	MOCK_PROVINCE	   = "@province"
	MOCK_CITY		   = "@city"
	MOCK_ZIP           = "@zip"
	MOCK_PICK          = "@pick"
	MOCK_ARRAY         = "@arr"

	STR_FEATURE_LOWER  = "lower"
	STR_FEATURE_UPPER  = "upper"
	STR_FEATURE_NUMBER = "number"
	STR_FEATURE_ALL    = "all"

	MAX_INT             = 1<<32 - 1
	MIN_INT             = 0 - MAX_INT - 1
	DEFAULT_DATE_FORMAT = "2006-01-02 15:04:05"
)

var StrMockFeatures map[string]int
var MOCK_MAP map[string]MockType
var mockManager *MockManager = &MockManager{}
var firstArr []string = []string{"赵","钱","孙","李","周","吴","郑","王","冯","陈","褚","卫","蒋","沈","韩","杨","朱","秦","尤","许",
								"何","吕","施","张","孔","曹","严","华","金","魏","陶","姜","戚","谢","邹","喻","柏","水","窦","章","云","苏","潘","葛","奚","范","彭","郎",
								"鲁","韦","昌","马","苗","凤","花","方","俞","任","袁","柳","酆","鲍","史","唐","费","廉","岑","薛","雷","贺","倪","汤","滕","殷",
								"罗","毕","郝","邬","安","常","乐","于","时","傅","皮","卞","齐","康","伍","余","元","卜","顾","孟","平","黄","和",
								"穆","萧","尹","姚","邵","湛","汪","祁","毛","禹","狄","米","贝","明","臧","计","伏","成","戴","谈","宋","茅","庞","熊","纪","舒",
								"屈","项","祝","董","梁","杜","阮","蓝","闵","席","季","麻","强","贾","路","娄","危","江","童","颜","郭","梅","盛","林","刁","钟",
								"徐","邱","骆","高","夏","蔡","田","樊","胡","凌","霍","虞","万","支","柯","昝","管","卢","莫","经","房","裘","缪","干","解","应",
								"宗","丁","宣","贲","邓","郁","单","杭","洪","包","诸","左","石","崔","吉","钮","龚","程","嵇","邢","滑","裴","陆","荣","翁","荀",
								"羊","于","惠","甄","曲","家","封","芮","羿","储","靳","汲","邴","糜","松","井","段","富","巫","乌","焦","巴","弓","牧","隗","山",
								"谷","车","侯","宓","蓬","全","郗","班","仰","秋","仲","伊","宫","宁","仇","栾","暴","甘","钭","厉","戎","祖","武","符","刘","景",
								"詹","束","龙","叶","幸","司","韶","郜","黎","蓟","溥","印","宿","白","怀","蒲","邰","从","鄂","索","咸","籍","赖","卓","蔺","屠",
								"蒙","池","乔","阴","郁","胥","能","苍","双","闻","莘","党","翟","谭","贡","劳","逄","姬","申","扶","堵","冉","宰","郦","雍","却",
								"璩","桑","桂","濮","牛","寿","通","边","扈","燕","冀","浦","尚","农","温","别","庄","晏","柴","瞿","阎","充","慕","连","茹","习",
								"宦","艾","鱼","容","向","古","易","慎","戈","廖","庾","终","暨","居","衡","步","都","耿","满","弘","匡","国","文","寇","广","禄",
								"阙","东","欧","殳","沃","利","蔚","越","夔","隆","师","巩","厍","聂","晁","勾","敖","融","冷","訾","辛","阚","那","简","饶","空",
								"曾","毋","沙","乜","养","鞠","须","丰","巢","关","蒯","相","查","后","荆","红","游","郏","竺","权","逯","盖","益","桓","公","仉",
								"督","岳","帅","缑","亢","况","郈","有","琴","归","海","晋","楚","闫","法","汝","鄢","涂","钦","商","牟","佘","佴","伯","赏","墨",
								"哈","谯","篁","年","爱","阳","佟","言","福","南","火","铁","迟","漆","官","冼","真","展","繁","檀","祭","密","敬","揭","舜","楼",
								"疏","冒","浑","挚","胶","随","高","皋","原","种","练","弥","仓","眭","蹇","覃","阿","门","恽","来","綦","召","仪","风","介","巨",
								"木","京","狐","郇","虎","枚","抗","达","杞","苌","折","麦","庆","过","竹","端","鲜","皇","亓","老","是","秘","畅","邝","还","宾",
								"闾","辜","纵","侴","万俟","司马","上官","欧阳","夏侯","诸葛","闻人","东方","赫连","皇甫","羊舌","尉迟","公羊","澹台","公冶","宗正",
								"濮阳","淳于","单于","太叔","申屠","公孙","仲孙","轩辕","令狐","钟离","宇文","长孙","慕容","鲜于","闾丘","司徒","司空","兀官","司寇",
								"南门","呼延","子车","颛孙","端木","巫马","公西","漆雕","车正","壤驷","公良","拓跋","夹谷","宰父","谷梁","段干","百里","东郭","微生",
								"梁丘","左丘","东门","西门","南宫","第五","公仪","公乘","太史","仲长","叔孙","屈突","尔朱","东乡","相里","胡母","司城","张廖","雍门",
								"毋丘","贺兰","綦毋","屋庐","独孤","南郭","北宫","王孙"}
var region []string = []string{"华北", "东北", "华东", "中南", "西南", "西北", "港澳台"}
var province map[string][]string = map[string][]string{"华北": {"北京市", "天津市", "河北省", "山西省", "内蒙古自治区"},
	"东北": {"辽宁省", "吉林省", "黑龙江省"},
	"华东": {"上海市", "江苏省", " 浙江省", "安徽省", "福建省", "江西省", "山东省"},
	"中南": {"河南省", "湖北省", "湖南省", "广东省", "广西壮族自治区", "海南省"},
	"西南": {"重庆市", "四川省", "贵州省", "云南省", "西藏自治区"},
	"西北": {"陕西省", "甘肃省", "青海省", "宁夏回族自治区", "新疆维吾尔自治区"},
	"港澳台": {"香港特别行政区", "澳门特别行政区", "台湾省"}}
var city map[string][]string = map[string][]string{"北京市": {"北京市"}, "天津市": {"天津市"},
	"河北省": {"石家庄", "保定市", "秦皇岛", "唐山市", "邯郸市", "邢台市", "沧州市", "承德市", "廊坊市", "衡水市", "张家口"},
	"山西省": {"太原市", "大同市", "阳泉市", "长治市", "临汾市", "晋中市", "运城市", "晋城市", "忻州市", "朔州市", "吕梁市 "},
	"内蒙古自治区": {"呼和浩特", "呼伦贝尔", "包头市", "赤峰市", "乌海市", "通辽市", "鄂尔多斯", "乌兰察布", "巴彦淖尔 "},
	"辽宁省": {"盘锦市", "鞍山市", "抚顺市", "本溪市", "铁岭市", "锦州市", "丹东市", "辽阳市", "葫芦岛", "阜新市", "朝阳市", "营口市"},
	"吉林省": {"吉林市", "通化市", "白城市", "四平市", "辽源市", "松原市", "白山市"},
	"黑龙江省": {"伊春市", "牡丹江", "大庆市", "鸡西市", "鹤岗市", "绥化市", "双鸭山", "七台河", "佳木斯", "黑河市", "齐齐哈尔市"},
	"上海市": {"上海市"},
	"江苏省": {"无锡市", "常州市", "扬州市", "徐州市", "苏州市", "连云港", "盐城市", "淮安市", "宿迁市", "镇江市", "南通市", "泰州市"},
	"浙江省": {"绍兴市", "温州市", "湖州市", "嘉兴市", "台州市", "金华市", "舟山市", "衢州市", "丽水市", "安徽省"},
	"安徽省": {"合肥市", "芜湖市", "亳州市", "马鞍山", "池州市", "淮南市", "淮北市", "蚌埠市", "巢湖市", "安庆市", "宿州市", "宣城市", "滁州市", "黄山市", "六安市", "阜阳市", "铜陵市"},
	"福建省": {"福州市", "泉州市", "漳州市", "南平市", "三明市", "龙岩市", "莆田市", "宁德市"},
	"江西省": {"南昌市", "赣州市", "景德镇", "九江市", "萍乡市", "新余市", "抚州市", "宜春市", "上饶市", "鹰潭市", "吉安市"},
	"山东省": {"潍坊市", "淄博市", "威海市", "枣庄市", "泰安市", "临沂市", "东营市", "济宁市", "烟台市", "菏泽市", "日照市", "德州市", "聊城市", "滨州市", "莱芜市"},
	"河南省": {"郑州市", "洛阳市", "焦作市", "商丘市", "信阳市", "新乡市", "安阳市", "开封市", "漯河市", "南阳市", "鹤壁市", "平顶山", "濮阳市", "许昌市", "周口市", "三门峡", "驻马店"},
	"湖北省": {"荆门市", "咸宁市", "襄樊市", "荆州市", "黄石市", "宜昌市", "随州市", "鄂州市", "孝感市", "黄冈市", "十堰市"},
	"湖南省": {"长沙市", "郴州市", "娄底市", "衡阳市", "株洲市", "湘潭市", "岳阳市", "常德市", "邵阳市", "益阳市", "永州市", "张家界", "怀化市"},
	"广东省": {"江门市", "佛山市", "汕头市", "湛江市", "韶关市", "中山市", "珠海市", "茂名市", "肇庆市", "阳江市", "惠州市", "潮州市", "揭阳市", "清远市", "东莞市", "汕尾市", "云浮市"},
	"广西壮族自治区": {"南宁市", "贺州市", "柳州市", "桂林市", "梧州市", "北海市", "玉林市", "钦州市", "百色市", "防城港", "贵港市", "河池市", "崇左市", "来宾市"},
	"海南省": {"海口市", "三亚市"},
	"重庆市": {"重庆市"},
	"四川省": {"乐山市", "雅安市", "广安市", "南充市", "自贡市", "泸州市", "内江市", "宜宾市", "广元市", "达州市", "资阳市", "绵阳市", "眉山市", "巴中市", "攀枝花", "遂宁市", "德阳市"},
	"贵州省": {"贵阳市", "安顺市", "遵义市", "六盘水"},
	"云南省": {"昆明市", "玉溪市", "大理市", "曲靖市", "昭通市", "保山市", "丽江市", "临沧市"},
	"西藏自治区": {"拉萨市", "阿里"},
	"陕西省": {"咸阳市", "榆林市", "宝鸡市", "铜川市", "渭南市", "汉中市", "安康市", "商洛市", "延安市"},
	"甘肃省": {"兰州市", "白银市", "武威市", "金昌市", "平凉市", "张掖市", "嘉峪关", "酒泉市", "庆阳市", "定西市", "陇南市", "天水市"},
	"青海省": {"西宁市"},
	"宁夏回族自治区": {"银川市", "固原市", "青铜峡市", "石嘴山市", "中卫市"},
	"新疆维吾尔自治区": {"乌鲁木齐", "克拉玛依市"},
	"香港特别行政区": {"香港岛", "九龙", "新界西"},
	"澳门特别行政区": {"澳门半岛", "澳门离岛"},
	"台湾省": {"基隆市", "台中市", "新竹市", "台南市", "嘉义市", "台北市", "高雄市"}}

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
	MOCK_MAP[MOCK_NAME] = &NameMock{firstMock: &FirstMock{}, lastMock: &LastMock{}}
	MOCK_MAP[MOCK_FIRST] = &FirstMock{}
	MOCK_MAP[MOCK_LAST] = &LastMock{}
	MOCK_MAP[MOCK_URL] = &UrlMock{}
	MOCK_MAP[MOCK_EMAIL] = &EmailMock{Arr: &[][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}}
	MOCK_MAP[MOCK_IP] = &IPMock{}
	MOCK_MAP[MOCK_REGION] = &RegionMock{}
	MOCK_MAP[MOCK_PROVINCE] = &ProvinceMock{}
	MOCK_MAP[MOCK_CITY] = &CityMock{}
	MOCK_MAP[MOCK_ADDRESS] = &AddressMock{regionMock: &RegionMock{}, provinceMock: &ProvinceMock{}, cityMock: &CityMock{}}
	MOCK_MAP[MOCK_ZIP] = &ZipMock{}
	MOCK_MAP[MOCK_PICK] = &PickMock{numMock: &NumMock{}}
	MOCK_MAP[MOCK_ARRAY] = &ArrayMock{numMock: &NumMock{}}
}

type MockManager struct{}

func GetMockManager() *MockManager {
	return mockManager
}

//func (mockManager *MockManager) Mock(str *string) (interface{}, error) {
//	if *str == "" {
//		return nil, errors.New("mock str is empty")
//	}
//	arr := REGEXP.FindStringSubmatch(*str)
//	length := len(arr)
//	if length < 2 {
//		return str, nil
//	}
//	mockType := arr[1]
//	if val, ok := MOCK_MAP[mockType]; ok {
//		if length > 2 {
//			arrtt := arr[2:]
//			return val.MockVal(&arrtt)
//		}
//		return val.MockVal(nil)
//	}
//	return str, nil
//}

func (mockManager *MockManager) IsSpecifiedAnnotation(str, mockStr string) bool {
	rs := []rune(str)
	index := strings.Index(str, MOCK_BRACKET_LEFT)
	preIndex := strings.Index(str, MOCK_PREFIX)
	prefixStr := string(rs[preIndex:index])
	return strings.Contains(prefixStr, mockStr)
}

func (mockManager *MockManager) MockDataFunc(str *string, f func(interface{}) (interface{}, error)) (interface{}, error) {
	if *str == "" {
		return nil, errors.New("mock str is empty")
	}
	index := strings.Index(*str, MOCK_BRACKET_LEFT)
	preIndex := strings.Index(*str, MOCK_PREFIX)
	lastIndex := mockManager.getFirstBracketRightIndex(str)
	rs := []rune(*str)
	prefix := string(rs[preIndex:index])
	resultPrefix := string(rs[:preIndex])
	resultLast := ""
	if lastIndex+1 < len(*str) {
		resultLast = string(rs[lastIndex+1:])
	}
	if val, ok := MOCK_MAP[prefix]; ok {
		var mockParams *[]string
		if index+1 == lastIndex {
			mockParams = nil
		} else {
			paramStr := string(rs[index+1:lastIndex])
			if paramStrTrim := strings.TrimSpace(paramStr); paramStrTrim == "" {
				mockParams = nil
			} else {
				mockParams = mockManager.getMockParams(&paramStrTrim)
			}
		}
		result, err := val.MockVal(mockParams)
		if err != nil {
			return result, err
		}
		if resultPrefix != "" || resultLast != "" {
			result, err = mockManager.warpResult(result, resultPrefix, resultLast)
			if err != nil {
				return nil, err
			}
		}
		if f == nil {
			return result, nil
		}
		return f(result)
	}
	return nil, errors.New("not support " + prefix)
}

func (mockManager *MockManager) MockData(str *string) (interface{}, error) {
	return mockManager.MockDataFunc(str, nil)
}

func (mockManager *MockManager) warpResult(res interface{}, start, end string) (interface{}, error) {
	var result string
	if val, ok := res.(string); ok {
		result = val
	}
	if val, ok := res.(int); ok {
		result = strconv.Itoa(val)
	}
	if result == "" {
		return nil, errors.New("just support string,int join")
	}
	if start != "" {
		startVal, err := getValue(start)
		if err != nil {
			return nil, err
		}
		result = startVal + result
	}
	if end != "" {
		endVal, err := getValue(end)
		if err != nil {
			return nil, err
		}
		result += endVal
	}
	return result, nil
}

func (mockManage *MockManager) getFirstBracketRightIndex(str *string) int {
	rs := []rune(*str)
	length := len(rs)
	pIndex := strings.Index(*str, MOCK_BRACKET_LEFT)
	time := 0
	index := pIndex
	for ; index < length; index++ {
		if rs[index] == '(' {
			time++
		} else if rs[index] == ')' {
			time--
		}
		if time == 0 {
			break
		}
	}
	if index >= length {
		return -1
	}
	return index
}

func (mockManager *MockManager) getMockParams(str *string) *[]string {
	arr := strings.Split(*str, MOCK_SPLIT)
	length := len(arr)
	result := make([]string, 0, length)
	tmpArr := make([]string, 0, length)
	inTime := 0
	for _, val := range arr {
		str := strings.TrimSpace(val)
		if tmp := strings.Count(str, MOCK_BRACKET_LEFT); tmp > 0 {
			inTime += tmp
		}
		if tmp := strings.Count(str, MOCK_BRACKET_RIGHT); tmp > 0 {
			inTime -= tmp
		}
		tmpArr = append(tmpArr, str)
		if inTime == 0 {
			result = append(result, strings.Join(tmpArr, MOCK_SPLIT))
			tmpArr = make([]string, 0, length)
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
	Arr *[][]int
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
	size := rand.Intn(max-min+1) + min
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
	ts := random.Intn(max-min+1) + min
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
	intpart = random.Intn(max-min+1) + min
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
		return fmt.Sprintf("%."+strconv.Itoa(dmax-dmin)+"f", float64(intpart)+random.Float64()), nil
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
		t = time.Unix(tNow-random.Int63n(tNow), 0)
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
			return nil, err
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
	firstMock 		*FirstMock
	lastMock 		*LastMock
}

func (nameMock NameMock) MockVal(params *[]string) (interface{}, error) {
	var first, last interface{}
	var err error
	if first, err = nameMock.firstMock.MockVal(nil); err != nil {
		return nil, err
	}
	if last, err = nameMock.lastMock.MockVal(nil); err != nil {
		return nil, err
	}
	return first.(string) + last.(string), err
}

type FirstMock struct {
	BaseMock
}

func (firstMock FirstMock) MockVal(params *[]string) (interface{}, error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return firstArr[random.Intn(len(firstArr))], nil
}

type LastMock struct {
	BaseMock
}

func (lastMock LastMock) MockVal(params *[]string) (interface{}, error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	times := random.Intn(2) + 1
	result := ""
	for i := 0; i < times; i++ {
		one := 0x4e00 + random.Int63n(0x9fa5 - 0x4e00 + 1)
		result += string(rune(one))
	}
	return result, nil
}

type UrlMock struct {
	BaseMock
}

func (urlMock UrlMock) MockVal(params *[]string) (interface{}, error) {
	str := "http://@str(5,8,lower).@str(2,3,lower)/@str(4,8,lower)"
	return GetMockManager().MockData(&str)
}

type EmailMock struct {
	BaseMock
	Arr *[][]int
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

type RegionMock struct {
	BaseMock
}

func (regionMock RegionMock) MockVal(params *[]string) (interface{}, error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return region[random.Intn(len(region))], nil
}

type ProvinceMock struct {
	BaseMock
}

func (provinceMock ProvinceMock) MockVal(params *[]string) (interface{}, error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	var value []string
	var ok bool
	if params == nil || len(*params) == 0 {
		for _, val := range province {
			value = val
			break
		}
	} else if value, ok = province[(*params)[0]]; !ok {
		return nil, errors.New("not find key: " + (*params)[0] + " in province")
	}
	return value[random.Intn(len(value))], nil
}

type CityMock struct {
	BaseMock
}

func (cityMock CityMock) MockVal(params *[]string) (interface{}, error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	var value []string
	var ok bool
	if params == nil || len(*params) == 0 {
		for _, val := range city {
			value = val
			break
		}
	} else if value, ok = city[(*params)[0]]; !ok {
		return nil, errors.New("not find key: " + (*params)[0] + " in city")
	}
	return value[random.Intn(len(value))], nil
}

type AddressMock struct {
	BaseMock
	regionMock 		*RegionMock
	provinceMock 	*ProvinceMock
	cityMock		*CityMock
}

func (addressMock AddressMock) MockVal(params *[]string) (interface{}, error) {
	regionStr, err := addressMock.regionMock.MockVal(nil)
	if err != nil {
		return nil, err
	}
	param := &[]string{regionStr.(string)}
	provinceStr, err := addressMock.provinceMock.MockVal(param)
	if err != nil {
		return nil, err
	}
	param = &[]string{provinceStr.(string)}
	cityStr, err := addressMock.cityMock.MockVal(param)
	if err != nil {
		return nil, err
	}
	return regionStr.(string) + provinceStr.(string) + cityStr.(string), nil
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
	numMock 		*NumMock
}

func (pickMock PickMock) MockVal(params *[]string) (interface{}, error) {
	paramArr := *params
	length := len(paramArr)
	if length == 0 {
		return nil, errors.New("no params")
	} else if length == 1 {
		rs := []rune(paramArr[0])
		last := len(rs) - 1
		return string(rs[1: last]), nil
	}
	rs := []rune(paramArr[0])
	paramArr[0] = strings.TrimSpace(string(rs[1:]))
	rs = []rune(paramArr[length -1])
	paramArr[length - 1] = strings.TrimSpace(string(rs[:len(rs) -1]))
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return paramArr[random.Intn(length)], nil
}

type IncrMock struct {
	BaseMock
}

type ArrayMock struct {
	BaseMock
	numMock		*NumMock
}

func (arrayMock ArrayMock) MockVal(params *[]string) (interface{}, error) {
	return arrayMock.numMock.MockVal(params)
}

func getValue(str string) (string, error) {
	val, err := get(str)
	if err != nil {
		return "", err
	}
	if val, ok := val.(string); ok {
		return val, nil
	}
	return "", errors.New("not string type")
}

func get(str string) (interface{}, error) {
	tmpStr := strings.TrimSpace(str)
	if strings.Contains(tmpStr, MOCK_PREFIX) {
		return mockManager.MockData(&tmpStr)
	}
	return tmpStr, nil
}

func getInt(str string) (int, error) {
	val, err := get(str)
	if err != nil {
		return 0, err
	}
	if val, ok := val.(int); ok {
		return val, nil
	}
	if val, ok := val.(string); ok {
		return strconv.Atoi(val)
	}
	return 0, errors.New("not int type")
}
