package user

import (
	"beebe/model/user"
	serivce "beebe/service"
)

type UserService interface {
	FindUserByUserId(userId *int64) *user.User;
	Login(loginUser *user.User) (*user.User, bool);
	RegisterUser(registerUser *user.User) bool
}

type UserServiceImpl struct {}


func (userService *UserServiceImpl) FindUserByUserId(userId *int64) *user.User{
	user := new(user.User)
	serivce.DB().First(user, userId)
	return user
}

func (userService *UserServiceImpl) Login(loginUser *user.User) (*user.User, bool) {
	user := new(user.User)
	serivce.DB().Where("name = ? and password = ?", loginUser.Name, loginUser.Password).Find(&user)
	if user.ID <= 0 {
		return nil, false
	}
	return user, true
}

func (userService *UserServiceImpl) RegisterUser(user *user.User) bool {
	return serivce.DB().Create(user).Error != nil
}


