package controller

import (
	macaron "gopkg.in/macaron.v1"
	"github.com/go-macaron/session"
	"github.com/martini-contrib/cors"
	"encoding/json"
	"beebe/model"
	"beebe/config"
	"strings"
	"fmt"
)

const (
	RESULT_KEY string = "result_data"
	ERROR_CODE_KEY	string = "error_code"
	ERROR_INFO_KEY string = "error_info"
)

var m *macaron.Macaron

func init() {
	m = macaron.New()
	setHandler()
	m.Use(macaron.Recovery())
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(macaron.Static("public"))
	sessionConfig()
}

func Macaron() *macaron.Macaron {
	return m
}


func setHandler() {
	if config.GetConfig().Web.IsCors {
		corsConfig := config.GetConfig().Cors
		m.Use(cors.Allow(&cors.Options{
			AllowOrigins:     strings.Split(corsConfig.AllowOrigin, ","),
			AllowMethods:     strings.Split(corsConfig.AllowMethods, ","),
			AllowHeaders:     strings.Split(corsConfig.AllowHeaders, ","),
			AllowCredentials: corsConfig.AllowCred,
		}))
	}
	m.Use(jsonResponse)
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
	ctx.Next()
	if errCode, ok := ctx.Data[ERROR_CODE_KEY]; ok {
		if value, ok := errCode.(model.ErrorCode); ok {
			restResult := model.ConvertRestResult(&value)
			if value.Code == model.SUCCESS.Code {
				if result, ok := ctx.Data[RESULT_KEY]; ok && result != nil {
					restResult.SetData(result)
				}
			}
			if err, ok := ctx.Data[ERROR_INFO_KEY]; ok {
				//log.Logger().Error("error", zap.Any("err", err))
				fmt.Printf("%v", err)
			}
			resultError, err := json.Marshal(*restResult)
			if err != nil {
				return ""
			}
			return string(resultError)
		}
	}
	return ""
}

func needLogin(ctx *macaron.Context, sess session.Store) {
	if user := getCurrentUser(sess); user == nil {
		ctx.Resp.Write(NoLoginResult)
	}
}

func noNeedLogin(ctx *macaron.Context, sess session.Store) {
	if user := getCurrentUser(sess); user != nil {
		ctx.Resp.Write(AlreadyLoginResult)
	}
}

func setResponse(ctx *macaron.Context, result interface{}, errCode *model.ErrorCode, err error) {
	if errCode != nil {
		ctx.Data[ERROR_CODE_KEY] = *errCode
		if errCode == model.SUCCESS && result != nil {
			ctx.Data[RESULT_KEY] = result
		}
		if err != nil {
			ctx.Data[ERROR_INFO_KEY] = err
		}
	}
}

func setSuccessResponse(ctx *macaron.Context, result interface{}) {
	setResponse(ctx, result, model.SUCCESS, nil)
}

func setFailResponse(ctx *macaron.Context, errCode *model.ErrorCode, err error) {
	setResponse(ctx, nil, errCode, err)
}

func setErrorResponse(ctx *macaron.Context, errCode *model.ErrorCode)  {
	setResponse(ctx, nil, errCode, nil)
}

