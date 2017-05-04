package controller

import (
	"gopkg.in/macaron.v1"
	"beebe/model"
	"beebe/service"
	"github.com/go-macaron/session"
	"github.com/go-macaron/binding"
)

type UserController struct {}
type UserLogin struct {
	UserName		string 			`json:"userName"`
	Password		string			`json:"password"`
}

func init() {
	userController := new(UserController)
	Macaron().Group("/user", func() {
		Macaron().Post("/login", binding.Bind(UserLogin{}), userController.login, jsonResponse)
		Macaron().Post("/logout", userController.logout, jsonResponse)
	})
}

func (userController *UserController) login(userLogin UserLogin, ctx *macaron.Context, sess session.Store) {
	userName := userLogin.UserName
	password := userLogin.Password
	if userName == "" || password == "" {
		ctx.Data[ERROR_CODE_KEY] = *model.PARAMETER_INVALID
		return
	}
	user, err := service.GetUserService().Login(&model.User{Name: userName, Password: password})
	if err != nil || user == nil {
		ctx.Data[ERROR_CODE_KEY] = *model.USERNAME_PASSWORD_ERROR
		ctx.Data[ERROR_INFO_KEY] = err
		return
	}
	if err := sess.Set(model.USER_SESSION_KEY, user); err != nil {
		ctx.Data[ERROR_CODE_KEY] = *model.SYSTEM_ERROR
		ctx.Data[ERROR_INFO_KEY] = err
		return
	}
	ctx.Data[RESULT_KEY] = &model.User{Name:user.Name, ID:user.ID, Email:user.Email}
}

func (userController *UserController) logout(ctx *macaron.Context, sess session.Store) {
	if err := sess.Delete(model.USER_SESSION_KEY); err != nil {
		ctx.Data[ERROR_CODE_KEY] = *model.SYSTEM_ERROR
		ctx.Data[ERROR_INFO_KEY] = err
		return
	}
	ctx.Data[RESULT_KEY] = *model.SUCCESS
}

func getCurrentUser(sess session.Store) *model.User {
	usertmp := sess.Get(model.USER_SESSION_KEY)
	if user, ok := usertmp.(model.User); ok {
		return *user
	}
	return nil
}
