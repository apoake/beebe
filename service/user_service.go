package service

import (
	"beebe/model"
	"github.com/pkg/errors"
)

var userService *UserServiceImpl = new(UserServiceImpl)

func GetUserService() *UserServiceImpl {
	return userService
}

type UserService interface {
	FindUserByUserId(userId *int64) *model.User;
	Login(loginUser *model.User) (*model.User, error);
	RegisterUser(registerUser *model.User) error
}

type UserServiceImpl struct {}


func (userService *UserServiceImpl) FindUserByUserId(userId *int64) *model.User{
	user := new(model.User)
	DB().First(user, userId)
	return user
}

func (userService *UserServiceImpl) Login(loginUser *model.User) (*model.User, error) {
	user := new(model.User)
	err := DB().Where("name = ? and password = ?", loginUser.Name, loginUser.Password).Find(user).Error
	return user, err
}

func (userService *UserServiceImpl) RegisterUser(user *model.User) error {
	return DB().Create(user).Error
}


type TeamService interface {
	Create(team *model.Team) error
	AddTeamUser(teamUser *model.TeamUser) error
	//MyTeam(userId *int64) (*[]model.Team, error)
	//MyJoiningTeam(userId *int64) (*[]model.Team, error)
}

type TeamServiceImpl struct {}

func (teamService *TeamServiceImpl) Create(team *model.Team) error {
	if team.UserId == 0 {
		return errors.New("param[team.userId] is empty")
	}
	return DB().Create(team).Error
}

func (teamService *TeamServiceImpl) AddTeamUser(teamUser *model.TeamUser) error {
	if teamUser == nil || teamUser.UserId == nil {

	}
}


