package service

import (
	"beebe/model"
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


