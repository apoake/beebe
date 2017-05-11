package controller

import (
	"gopkg.in/macaron.v1"
	"beebe/model"
	"beebe/service"
	"github.com/go-macaron/session"
	"github.com/go-macaron/binding"
	"encoding/json"
	"beebe/utils"
)

var NoLoginResult []byte
var AlreadyLoginResult []byte

type UserController struct {}
type UserDto struct {
	Account			string 				`json:"userName"`
	Password		string				`json:"password"`
	Email			string				`json:"email"`
}

func init() {
	NoLoginResult, _ = json.Marshal( model.ConvertRestResult(model.USER_NO_LOGIN))
	AlreadyLoginResult, _ = json.Marshal(model.ConvertRestResult(model.USER_ALREADY_LOGIN))
	userController := new(UserController)
	Macaron().Group("/user", func() {
		Macaron().Post("/register", noNeedLogin, binding.Bind(UserDto{}), userController.register)
		Macaron().Post("/login", noNeedLogin, binding.Bind(UserDto{}), userController.login)
		Macaron().Post("/logout", needLogin, userController.logout)
	})
}

func (userController *UserController) register(user UserDto, ctx *macaron.Context) {
	if user.Account == "" || user.Password == "" {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	if service.GetUserService().CheckUserByAccount(&user.Account) {
		setErrorResponse(ctx, model.USER_ACCOUNT_EXIST)
		return
	}
	err := service.GetUserService().RegisterUser(&model.User{Account: user.Account, Password: utils.SHA(user.Password), Email: user.Email})
	if err != nil {
		setFailResponse(ctx, model.USER_REGISTER_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

func (userController *UserController) login(userLogin UserDto, ctx *macaron.Context, sess session.Store) {
	userName := userLogin.Account
	password := userLogin.Password
	if userName == "" || password == "" {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	user, err := service.GetUserService().Login(&model.User{Account: userName, Password: utils.SHA(password)})
	if err != nil || user == nil {
		setFailResponse(ctx, model.USERNAME_PASSWORD_ERROR, err)
		return
	}
	if err := sess.Set(model.USER_SESSION_KEY, user); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, &model.User{Name:user.Name, ID:user.ID, Email:user.Email, Account:user.Account})
}

func (userController *UserController) logout(ctx *macaron.Context, sess session.Store) {
	if err := sess.Delete(model.USER_SESSION_KEY); err != nil {
		setResponse(ctx, nil, model.SYSTEM_ERROR, err)
		return
	}
	setErrorResponse(ctx, model.SUCCESS)
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

func getCurrentUser(sess session.Store) *model.User {
	usertmp := sess.Get(model.USER_SESSION_KEY)
	if user, ok := usertmp.(*model.User); ok {
		return user
	}
	return nil
}

func getCurrentUserId(sess session.Store) *int64 {
	user := getCurrentUser(sess)
	if user != nil {
		return nil
	}
	return &(user.ID)
}
