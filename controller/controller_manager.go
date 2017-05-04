package controller

import (
	macaron "gopkg.in/macaron.v1"
	"github.com/go-macaron/session"
	"encoding/json"
	"github.com/eleme/taco/log"
	"beebe/model"
)

const (
	RESULT_KEY string = "result_data"
	ERROR_CODE_KEY	string = "error_code"
	ERROR_INFO_KEY string = "error_info"
)

var m *macaron.Macaron

func init() {
	m = macaron.Classic()
	m.Use(macaron.Recovery())
	sessionConfig()
}

func Macaron() *macaron.Macaron {
	return m
}

func sessionConfig() {
	m.Use(session.Sessioner(session.Options{
		// 提供器的名称，默认为 "memory"
		Provider:       "memory",
		// 提供器的配置，根据提供器而不同
		ProviderConfig: "",
		// 用于存放会话 ID 的 Cookie 名称，默认为 "MacaronSession"
		CookieName:     "gsession",
		// Cookie 储存路径，默认为 "/"
		CookiePath:     "/",
		// GC 执行时间间隔，默认为 3600 秒
		Gclifetime:     3600,
		// 最大生存时间，默认和 GC 执行时间间隔相同
		Maxlifetime:    3600,
		// 仅限使用 HTTPS，默认为 false
		Secure:         false,
		// Cookie 生存时间，默认为 0 秒
		CookieLifeTime: 0,
		// Cookie 储存域名，默认为空
		Domain:         "",
		// 会话 ID 长度，默认为 16 位
		IDLength:       16,
		// 配置分区名称，默认为 "session"
		Section:        "session",
	}))
}

func jsonResponse(ctx *macaron.Context) string {
	if errCode, ok := ctx.Data[ERROR_CODE_KEY]; ok {
		if value, ok := errCode.(model.ErrorCode); ok {
			resultError, err := json.Marshal(*model.ConvertRestResult(&value))
			if err != nil {
				return ""
			}
			// TODO logger error
			return string(resultError)
		}
		return ""
	}
	result, ok := ctx.Data[RESULT_KEY]
	if !ok {
		resultError, err := json.Marshal(*model.ConvertRestResult(model.SYSTEM_ERROR))
		if err != nil {
			return ""
		}
		return string(resultError)
	}
	resultData, err := json.Marshal(*(model.ConvertRestResult(model.SUCCESS).SetData(result)))
	if err != nil {
		// TODO logger error
		log.Logger().Error("jsonResponse error")
	}
	return string(resultData)
}