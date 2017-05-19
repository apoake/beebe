package controller

import (
	macaron "gopkg.in/macaron.v1"
	"github.com/go-macaron/session"
	"github.com/martini-contrib/cors"
	"beebe/config"
	"strings"
	"time"
	"fmt"
	"beebe/log"
	"go.uber.org/zap"
	"net/http"
	"beebe/utils"
)

var m *macaron.Macaron
var LogTimeFormat string = "2006-01-02 15:04:05"

func init() {
	m = macaron.New()
	setHandler()
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
	//m.Use(macaron.Logger())
	m.Use(Logger())
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

func Logger() macaron.Handler {
	return func(ctx *macaron.Context) {
		start := time.Now()
		requestId := utils.GetGuid()
		log.Mlog.Info(fmt.Sprintf("Macaron %s:[%s] Started %s %s for %s", start.Format(LogTimeFormat), requestId, ctx.Req.Method, ctx.Req.RequestURI, ctx.RemoteAddr()),
			zap.String("requestId", requestId),
			zap.String("startTime", start.Format(LogTimeFormat)),
			zap.String("requestMethod", ctx.Req.Method),
			zap.String("requestUri", ctx.Req.RequestURI),
			zap.String("remoteAddr", ctx.RemoteAddr()))
		log.Mlog.Sync()
		ctx.Next()
		rw := ctx.Resp.(macaron.ResponseWriter)
		end := time.Now()
		status := rw.Status()
		statusName := http.StatusText(status)
		dur := time.Since(start).String()

		log.Mlog.Info(fmt.Sprintf("Macaron %s:[%s] Completed %s %d %s in %s", end.Format(LogTimeFormat), requestId, ctx.Req.RequestURI, rw.Status(), statusName, dur),
			zap.String("requestId", requestId),
			zap.String("endTime", end.Format(LogTimeFormat)),
			zap.String("requestUri", ctx.Req.RequestURI),
			zap.Int("status", status),
			zap.String("statusName", statusName),
			zap.String("useTime", dur))
		log.Mlog.Sync()
	}
}