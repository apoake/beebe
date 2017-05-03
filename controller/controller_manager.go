package controller

import (
	macaron "gopkg.in/macaron.v1"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/session"
)

var m *macaron.Macaron

func init() {
	m = macaron.Classic()
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