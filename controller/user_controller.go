package controller

import (
	"gopkg.in/macaron.v1"
	"beebe/model"
	"beebe/service"
	"github.com/go-macaron/session"
	"github.com/go-macaron/binding"
	"encoding/json"
)

var NoLoginResult []byte
var alreadyLoginResult []byte
type UserController struct {}
type UserLogin struct {
	UserName		string 			`json:"userName"`
	Password		string			`json:"password"`
}

func init() {
	noLoginResult := model.ConvertRestResult(model.USER_NO_LOGIN)
	NoLoginResult, _ = json.Marshal(noLoginResult)
	userController := new(UserController)
	Macaron().Group("/user", func() {
		Macaron().Post("/login", noNeedLogin, binding.Bind(UserLogin{}), userController.login, jsonResponse)
		Macaron().Post("/logout", needLogin, userController.logout, jsonResponse)
	})
}

func (userController *UserController) login(userLogin UserLogin, ctx *macaron.Context, sess session.Store) {
	userName := userLogin.UserName
	password := userLogin.Password
	if userName == "" || password == "" {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	user, err := service.GetUserService().Login(&model.User{Name: userName, Password: password})
	if err != nil || user == nil {
		setFailResponse(ctx, model.USERNAME_PASSWORD_ERROR, err)
		return
	}
	if err := sess.Set(model.USER_SESSION_KEY, user); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, &model.User{Name:user.Name, ID:user.ID, Email:user.Email})
}

func (userController *UserController) logout(ctx *macaron.Context, sess session.Store) {
	if err := sess.Delete(model.USER_SESSION_KEY); err != nil {
		setResponse(ctx, nil, model.SYSTEM_ERROR, err)
		return
	}
	setErrorResponse(ctx, model.SUCCESS)
}

func needLogin(ctx *macaron.Context, sess *session.Store) {
	if user := getCurrentUser(sess); user == nil {
		ctx.Resp.Write(NoLoginResult)
	}
}

func noNeedLogin(ctx *macaron.Context, sess, store *session.Store) {
	if user := getCurrentUser(sess); user != nil {
		ctx.Resp.Write(NoLoginResult)
	}
}

func getCurrentUser(sess session.Store) *model.User {
	usertmp := sess.Get(model.USER_SESSION_KEY)
	if user, ok := usertmp.(model.User); ok {
		return &user
	}
	return nil
}

func getCurrentUserId(sess session.Store) *int64 {
	user := getCurrentUser(sess)
	if user != nil {
		return nil
	}
	return &user.ID
}
