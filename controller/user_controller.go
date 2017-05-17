package controller

import (
	"gopkg.in/macaron.v1"
	"beebe/model"
	"beebe/service"
	"github.com/go-macaron/session"
	"github.com/go-macaron/binding"
	"beebe/utils"
	"errors"
)


type UserController struct{}

func init() {
	userController := new(UserController)
	Macaron().Group("/user", func() {
		Macaron().Post("/", needLogin, userController.user)
		Macaron().Post("/register", noNeedLogin, binding.Bind(UserRegister{}), userController.register)
		Macaron().Post("/login", noNeedLogin, binding.Bind(UserLogin{}), userController.login)
		Macaron().Post("/search", needLogin, binding.Bind(UserSearch{}), userController.search)
		Macaron().Post("/update", needLogin, binding.Bind(UserUpdate{}), userController.update)
		Macaron().Post("/changepassword", needLogin, binding.Bind(UserPassword{}), userController.changePassword)
		Macaron().Post("/logout", needLogin, userController.logout)
	})
	Macaron().Group("/team", func() {
		Macaron().Post("/create", binding.Bind(TeamAdd{}), userController.createTeam)
		Macaron().Post("/update", binding.Bind(TeamUpdate{}), userController.updateTeam)
		Macaron().Post("/info", binding.Bind(Id{}), userController.teamInfo)
		Macaron().Post("/mine", userController.myTeam)
		Macaron().Post("/join", userController.myJoinTeam)
		Macaron().Post("/adduser", binding.Bind(TeamUserDto{}), userController.addTeamUser)
		Macaron().Post("/removeuser", binding.Bind(TeamUserDto{}), userController.removeTeamUser)
	}, needLogin)
}

/**
	获取当前用户信息
 */
func (userController *UserController) user(ctx *macaron.Context, sess session.Store) {
	setSuccessResponse(ctx, getCurrentUser(sess))
}

/**
	用户注册
 */
func (userController *UserController) register(user UserRegister, ctx *macaron.Context) {
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

/**
	用户登入
 */
func (userController *UserController) login(userLogin UserLogin, ctx *macaron.Context, sess session.Store) {
	user, ok := service.GetUserService().Login(&model.User{Account: userLogin.Account, Password: utils.SHA(userLogin.Password)})
	if !ok {
		setErrorResponse(ctx, model.USERNAME_PASSWORD_ERROR)
		return
	}
	if err := sess.Set(model.USER_SESSION_KEY, model.User{ID: user.ID, Name: user.Name, ImgUrl: user.ImgUrl, Account: user.Account, Email: user.Email}); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

/**
	账户搜索
 */
func (userController *UserController) search(userSearch UserSearch, ctx *macaron.Context) {
	var limit int64 = 5
	users, err := service.GetUserService().SearchUserByAccount(&userSearch.Account, &limit)
	if err != nil {
		setErrorResponse(ctx, model.SYSTEM_ERROR)
		return
	}
	setSuccessResponse(ctx, users)
}

/**
	用户更新
 */
func (userController *UserController) update(userUpdate UserUpdate, ctx *macaron.Context, sess session.Store) {
	userSession := getCurrentUser(sess)
	user := &model.User{Name: userUpdate.Name, ImgUrl: userUpdate.ImgUrl, ID: userSession.ID, Email: userUpdate.Email}
	if err := service.GetUserService().UpdateUser(user); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	if err := sess.Set(model.USER_SESSION_KEY, model.User{ID: userSession.ID, Name: user.Name, ImgUrl: user.ImgUrl, Account: userSession.Account, Email: user.Email}); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

/**
	修改密码
 */
func (userController *UserController) changePassword(userPassword UserPassword, ctx *macaron.Context, sess session.Store) {
	userId := getCurrentUserId(sess)
	user, _ := service.GetUserService().FindUserByUserId(&userId)
	if pa := utils.SHA(userPassword.Opassword); pa != user.Password {
		setErrorResponse(ctx, model.USER_PASSWORD_DISAGREE)
		return
	}
	if err := service.GetUserService().ChangePassword(&model.User{ID: userId, Password: utils.SHA(userPassword.Password)}); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

/**
	用户登出
 */
func (userController *UserController) logout(ctx *macaron.Context, sess session.Store) {
	if err := sess.Delete(model.USER_SESSION_KEY); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setErrorResponse(ctx, model.SUCCESS)
}

/**
	创建团队
 */
func (userController *UserController) createTeam(teamAdd TeamAdd, ctx *macaron.Context, sess session.Store) {
	userId := getCurrentUserId(sess)
	err := service.GetTeamService().Create(&model.Team{Name: teamAdd.Name, UserId: userId, Remark: teamAdd.Remark, LogoUrl: teamAdd.LogoUrl})
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

func (userController *UserController) updateTeam(teamUpdate TeamUpdate, ctx *macaron.Context, sess session.Store) {
	teamDB, ok := service.GetTeamService().Get(&teamUpdate.ID)
	if !ok {
		setErrorResponse(ctx, model.TEAM_NO_EXIST)
		return
	}
	if teamDB.UserId != getCurrentUserId(sess) {
		setErrorResponse(ctx, model.TEAM_NOT_SELF)
		return
	}
	if err := service.GetTeamService().Update(&model.Team{ID: teamUpdate.ID, Name: teamUpdate.Name, Remark:teamUpdate.Remark, LogoUrl: teamUpdate.LogoUrl}); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

func (userController *UserController) teamInfo(id Id, ctx *macaron.Context, sess session.Store) {
	if !service.GetTeamService().HasTeamRight(&model.TeamUser{UserId: getCurrentUserId(sess), TeamId: id.ID}) {
		setErrorResponse(ctx, model.USER_NO_RIGHT)
		return
	}
	result, err := service.GetTeamService().Info(&model.Team{ID: id.ID})
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

func (userController *UserController) addTeamUser(teamUserDto TeamUserDto, ctx *macaron.Context, sess session.Store) {
	if teamUserDto.RoleId == 0 {
		teamUserDto.RoleId = model.ROLE_MEMBER.ID
	}
	param := &model.TeamUser{UserId: getCurrentUserId(sess), TeamId: teamUserDto.TeamId}
	if !service.GetTeamService().HasTeamRight(param) {
		setErrorResponse(ctx, model.USER_NO_RIGHT)
		return
	}
	if _, ok := service.GetTeamService().GetTeamUserByUserIdAndTeamId(param); ok {
		setErrorResponse(ctx, model.TEAM_ALREADY_JOIN)
		return
	}
	err := service.GetTeamService().AddTeamUser(&model.TeamUser{TeamId: teamUserDto.TeamId, RoleId: teamUserDto.RoleId, UserId: teamUserDto.UserId})
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

func (userController *UserController) removeTeamUser(teamUserDto TeamUserDto, ctx *macaron.Context, sess session.Store) {
	if teamUserDto.TeamId == 0 {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	if !service.GetTeamService().HasTeamRight(&model.TeamUser{UserId: getCurrentUserId(sess), TeamId: teamUserDto.TeamId}) {
		setErrorResponse(ctx, model.USER_NO_RIGHT)
		return
	}
	err := service.GetTeamService().RemoveTeamUser(&model.TeamUser{TeamId: teamUserDto.TeamId, UserId: teamUserDto.UserId})
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