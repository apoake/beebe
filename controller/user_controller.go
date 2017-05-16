package controller

import (
	"gopkg.in/macaron.v1"
	"beebe/model"
	"beebe/service"
	"github.com/go-macaron/session"
	"github.com/go-macaron/binding"
	"encoding/json"
	"beebe/utils"
	"errors"
)

var NoLoginResult []byte
var AlreadyLoginResult []byte

type UserController struct{}
type UserDto struct {
	Account  	string			`json:"userName"`
	Opassword 	string			`json:"opassword"`
	Password 	string          `json:"password"`
	Email    	string          `json:"email"`
	ImgUrl  	string			`json:"imgUrl"`
	Name 		string			`json:"nickName"`
}


func init() {
	NoLoginResult, _ = json.Marshal(model.ConvertRestResult(model.USER_NO_LOGIN))
	AlreadyLoginResult, _ = json.Marshal(model.ConvertRestResult(model.USER_ALREADY_LOGIN))
	userController := new(UserController)
	Macaron().Group("/user", func() {
		Macaron().Post("/", needLogin, userController.user)
		Macaron().Post("/register", noNeedLogin, binding.Bind(UserDto{}), userController.register)
		Macaron().Post("/login", noNeedLogin, binding.Bind(UserDto{}), userController.login)
		Macaron().Post("/search", needLogin, binding.Bind(UserDto{}), userController.search)
		Macaron().Post("/update", needLogin, binding.Bind(UserDto{}), userController.update)
		Macaron().Post("/changepassword", needLogin, binding.Bind(UserDto{}), userController.changePassword)
		Macaron().Post("/logout", needLogin, userController.logout)
	})
	Macaron().Group("/team", func() {
		Macaron().Post("/create", binding.Bind(model.Team{}), userController.createTeam)
		Macaron().Post("/update", binding.Bind(model.Team{}), userController.updateTeam)
		Macaron().Post("/info", binding.Bind(model.Team{}), userController.teamInfo)
		Macaron().Post("/mine", userController.myTeam)
		Macaron().Post("/join", userController.myJoinTeam)
		Macaron().Post("/adduser", binding.Bind(model.TeamUser{}), userController.addTeamUser)
		Macaron().Post("/removeuser", binding.Bind(model.TeamUser{}), userController.removeTeamUser)
	}, needLogin)
}

func (userController *UserController) user(ctx *macaron.Context, sess session.Store) {
	setSuccessResponse(ctx, getCurrentUser(sess))
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
	err := service.GetUserService().RegisterUser(&model.User{Account: user.Account, Password: utils.SHA(user.Password), Email: user.Email, Name: user.Account})
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
	user, ok := service.GetUserService().Login(&model.User{Account: userName, Password: utils.SHA(password)})
	if ok {
		setErrorResponse(ctx, model.USERNAME_PASSWORD_ERROR)
		return
	}
	if err := sess.Set(model.USER_SESSION_KEY, *user); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, model.User{Name: user.Name, ID: user.ID, Email: user.Email, Account: user.Account})
}

func (userController *UserController) search(userDto UserDto, ctx *macaron.Context) {
	accountName := userDto.Account
	if accountName == "" {
		setSuccessResponse(ctx, nil)
		return
	}
	var limit int64 = 5
	users, err := service.GetUserService().SearchUserByAccount(&accountName, &limit)
	if err != nil {
		setErrorResponse(ctx, model.SYSTEM_ERROR)
		return
	}
	setSuccessResponse(ctx, users)
}

func (userController *UserController) update(userDto UserDto, ctx *macaron.Context, sess session.Store) {
	userId := getCurrentUserId(sess)
	if userDto.Email == "" || userDto.ImgUrl == "" ||
		userDto.Name == "" {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	user := &model.User{Name: userDto.Name, ImgUrl: userDto.ImgUrl, ID: userId, Email: userDto.Email}
	if err := service.GetUserService().UpdateUser(user); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

func (userController *UserController) changePassword(userDto *UserDto, ctx *macaron.Context, sess session.Store) {
	if userDto.Password == "" || userDto.Opassword == "" {
		setErrorResponse(ctx, model.SYSTEM_ERROR)
		return
	}
	userId := getCurrentUserId(sess)
	user, _ := service.GetUserService().FindUserByUserId(&userId)
	if pa := utils.SHA(userDto.Password); pa != user.Password {
		setErrorResponse(ctx, model.USER_PASSWORD_DISAGREE)
		return
	}
	if err := service.GetUserService().ChangePassword(&model.User{Password: userDto.Password}); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

func (userController *UserController) logout(ctx *macaron.Context, sess session.Store) {
	if err := sess.Delete(model.USER_SESSION_KEY); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setErrorResponse(ctx, model.SUCCESS)
}

func (userController *UserController) createTeam(team model.Team, ctx *macaron.Context, sess session.Store) {
	if team.Name == "" {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	userId := getCurrentUserId(sess)
	err := service.GetTeamService().Create(&model.Team{Name: team.Name, UserId: userId, Remark: team.Name})
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

func (userController *UserController) updateTeam(team model.Team, ctx *macaron.Context, sess session.Store) {
	if team.Name == "" || team.ID == 0 {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	teamDB, ok := service.GetTeamService().Get(&team.ID)
	if !ok {
		setErrorResponse(ctx, model.TEAM_NO_EXIST)
		return
	}
	if teamDB.UserId != getCurrentUserId(sess) {
		setErrorResponse(ctx, model.TEAM_NOT_SELF)
		return
	}
	if err := service.GetTeamService().Update(&team); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

func (userController *UserController) teamInfo(team model.Team, ctx *macaron.Context, sess session.Store) {
	if team.ID == 0 {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	if !service.GetTeamService().HasTeamRight(&model.TeamUser{UserId: getCurrentUserId(sess), TeamId: team.ID}) {
		setErrorResponse(ctx, model.USER_NO_RIGHT)
		return
	}
	result, err := service.GetTeamService().Info(&team)
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, result)
}


func (userController *UserController) myTeam(ctx *macaron.Context, sess session.Store) {
	userId := getCurrentUserId(sess)
	teams, err := service.GetTeamService().MyTeam(&userId)
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, teams)
}

func (userController *UserController) myJoinTeam(ctx *macaron.Context, sess session.Store) {
	userId := getCurrentUserId(sess)
	teams, err := service.GetTeamService().MyJoinTeam(&userId)
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, teams)
}

func (userController *UserController) addTeamUser(teamUser model.TeamUser, ctx *macaron.Context, sess session.Store) {
	if teamUser.TeamId == 0 || teamUser.UserId == 0 {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	if teamUser.RoleId == 0 {
		teamUser.RoleId = model.ROLE_MEMBER.ID
	}
	if !service.GetTeamService().HasTeamRight(&model.TeamUser{UserId: getCurrentUserId(sess), TeamId: teamUser.TeamId}) {
		setErrorResponse(ctx, model.USER_NO_RIGHT)
		return
	}
	err := service.GetTeamService().AddTeamUser(&model.TeamUser{TeamId: teamUser.TeamId, RoleId: teamUser.RoleId, UserId: teamUser.UserId})
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

func (userController *UserController) removeTeamUser(teamUser model.TeamUser, ctx *macaron.Context, sess session.Store) {
	if teamUser.TeamId == 0 {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	if !service.GetTeamService().HasTeamRight(&model.TeamUser{UserId: getCurrentUserId(sess), TeamId: teamUser.TeamId}) {
		setErrorResponse(ctx, model.USER_NO_RIGHT)
		return
	}
	err := service.GetTeamService().RemoveTeamUser(&model.TeamUser{TeamId: teamUser.TeamId, UserId: teamUser.UserId})
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

func getCurrentUser(sess session.Store) *model.User {
	usertmp := sess.Get(model.USER_SESSION_KEY)
	if user, ok := usertmp.(model.User); ok {
		return &user
	}
	return nil
}

func getCurrentUserId(sess session.Store) int64 {
	user := getCurrentUser(sess)
	if user == nil {
		return 0
	}
	return user.ID
}

func getUserId(ctx *macaron.Context, sess session.Store) (int64, bool) {
	user := getCurrentUser(sess)
	if user == nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, errors.New("not find user in session"))
		return 0, false
	}
	return user.ID, true
}