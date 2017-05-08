package model

import (
	"time"
)

const (
	USER_SESSION_KEY = "user"
)

type User struct {
	Model
	ID				int64		`gorm:"primary_key" json:"id"`
	Account			string		`json:"account"`
	Password		string		`json:"password" form:"password"`
	Name 			string		`json:"name" form:"name"`
	Email 			string		`json:"email"`
	IsLockOut		uint		`gorm:"column:is_locked_out" json:"isLockOut" gorm:"default:0"`
	LastLoginDate	time.Time	`gorm:"column:last_login_date"`
}

type UserVo struct {
	Vo
	ID				int64		`json:"id"`
	Name 			string		`json:"name"`
	Email 			string		`json:"email"`
}

func (User) TableName() string {
	return "user"
}


type Role struct {
	Model
	ID 			int64	`gorm:"primary_key" json:"id"`
	RoleName	string	`grom:"column:"role_name" json:"roleName"`
}

func (Role) TableName() string {
	return "role"
}

type TeamUser struct {
	Model
	UserId 			int64		`grom:"column:"user_id" json:"userId"`
	TeamId			int64		`grom:"column:"team_id" json:"teamId"`
	RoleId			int64		`grom:"column:"role_Id" json:roleId"`
	ProjectId 		int64		`grom:"column:"project_id" json:projectId"`
}

func (TeamUser) TableName() string {
	return "team_user"
}

type UserSetting struct {
	Model
	UserId 		int64		`grom:"column:"user_id" json:"userId"`
	key 		string		`grom:"column:"key" json:"key"`
	val 		string		`grom:"column:"val" json:"val"`
}

func (UserSetting) TableName() string {
	return "user_settings"
}

type Team struct {
	Model
	ID 			int64		`gorm:"primary_key" json:"id"`
	Name 		string		`grom:"column:"name" json:"name"`
	Remark		string		`grom:"column:"remark" json:"remark"`
	LogoUrl		string		`grom:"column:"logo_url" json:"logoUrl"`
	UserId		int64		`grom:"column:"user_id" json:"userId"`
}

func (Team) TableName() string {
	return "team"
}

type WorkSpace struct {
	Model
	ID 			int64		`gorm:"primary_key" json:"id"`
	UserId		int64		`grom:"column:"user_id" json:"userId"`
	ProjectId	int64		`grom:"column:"project_id" json:"projectId"`
}

func (WorkSpace) TableName() string {
	return "workspace"
}
