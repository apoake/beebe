package service

import (
	"beebe/model"
	"github.com/pkg/errors"
)

var userService *UserServiceImpl = new(UserServiceImpl)
var teamService *TeamServiceImpl = new(TeamServiceImpl)

func GetUserService() *UserServiceImpl {
	return userService
}

func GetTeamService() *TeamServiceImpl {
	return teamService
}

type UserService interface {
	FindUserByUserId(userId *int64) *model.User;
	CheckUserByAccount(account *string) bool;
	Login(loginUser *model.User) (*model.User, error);
	RegisterUser(registerUser *model.User) error
	SearchUserByAccount(account *string) (*model.User, error)
}

type UserServiceImpl struct {}


func (userService *UserServiceImpl) FindUserByUserId(userId *int64) *model.User{
	user := new(model.User)
	DB().First(user, userId)
	return user
}

func (userService *UserServiceImpl) CheckUserByAccount(account *string) bool {
	user := new(model.User)
	DB().Where("account = ?", account).First(user)
	return user.ID > 0
}

func (userService *UserServiceImpl) Login(loginUser *model.User) (*model.User, error) {
	user := new(model.User)
	err := DB().Where("account = ? and password = ?", loginUser.Account, loginUser.Password).Find(user).Error
	return user, err
}

func (userService *UserServiceImpl) RegisterUser(user *model.User) error {
	return DB().Create(user).Error
}

func (userService *UserServiceImpl) SearchUserByAccount(account *string) (*model.User, error) {
	user := new(model.User)
	err := DB().Where("account = ?", account).First(user).Error
	return user, err
}


type TeamService interface {
	Get(teamId *int64) *model.Team
	Create(team *model.Team) error
	Update(team *model.Team) error
	QuitTeam(team *model.Team) error
	Transform(team *model.Team, userId *int64) error
	ChangeRole(teamUser *model.TeamUser) error
	MyTeam(userId *int64) (*[]model.Team, error)
	MyJoinTeam(userId *int64) (*[]model.Team, error)
	AddTeamUser(teamUser *model.TeamUser) error
	RemoveTeamUser(teamUser *model.TeamUser) error
}

type TeamServiceImpl struct {}

func (teamService *TeamServiceImpl) Get(teamId *int64) *model.Team {
	team := &model.Team{ID: *teamId}
	DB().First(team)
	return team
}

func (teamService *TeamServiceImpl) Create(team *model.Team) (err error) {
	if team.UserId == 0 {
		return errors.New("param[team.userId] is empty")
	}
	tx := DB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Create(team).Error; err != nil {
		return err
	}
	if err = tx.Create(&model.TeamUser{UserId: team.UserId, TeamId: team.ID, RoleId: model.ROLE_OWNER.ID}).Error; err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (teamService *TeamServiceImpl) Update(team *model.Team) error {
	if team.ID == 0 {
		return errors.New("param[team.ID] is empty")
	}
	return DB().Model(&model.Team{ID: team.ID}).Updates(map[string]interface{}{"name": team.Name, "remark": team.Remark}).Error
}

func (teamService *TeamServiceImpl) QuitTeam(team *model.Team) (err error) {
	if team.UserId == 0 {
		return errors.New("param[team.userId] is empty")
	}
	dbTeam := &model.Team{ID: team.ID}
	if err = DB().First(dbTeam).Error; err != nil {
		return
	}
	if team.UserId == dbTeam.UserId {
		err = errors.New("role owner; can not to delete")
		return 
	}
	tx := DB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Where("user_id = ? and team_id = ?", team.UserId, team.ID).Delete(&model.ProjectUserMapping{}).Error; err != nil {
		return
	}
	if err = tx.Where("user_id = ? and team_id = ?", team.UserId, team.ID).Delete(&model.TeamUser{}).Error; err != nil {
		return 
	}
	return nil
}

func (teamService *TeamServiceImpl) Transform(team *model.Team, userId *int64) (err error) {
	if team.UserId == 0 || team.ID == 0 || userId == nil {
		return errors.New("params error")
	}
	dbTeam := &model.Team{ID: team.ID}
	if err = DB().Find(dbTeam).Error; err != nil {
		return
	}
	if dbTeam.UserId != team.UserId {
		return errors.New("user is not the team owner")
	}
	roleUser := &model.TeamUser{}
	if err = DB().Where("user_id = ? and team_id = ?").First(roleUser).Error; err != nil {
		return err
	}
	if roleUser.ID < 0 {
		return errors.New("not find")
	}
	tx := db.Begin()
	dbTeam.UserId = *userId
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Save(dbTeam).Error; err != nil {
		return
	}
	if err = tx.Model(&model.TeamUser{}).
		Where("user_id = ? and team_id = ?", team.UserId, team.ID).
		Update("role_id", model.ROLE_MEMBER).Error; err != nil {
		return
	}
	if err = tx.Model(&model.TeamUser{}).
		Where("user_id = ? and team_id = ?", *userId, team.ID).
		Update("role_id", model.ROLE_OWNER).Error; err != nil {
		return
	}
	tx.Commit()
	return nil
}

func (teamService *TeamServiceImpl) ChangeRole(teamUser *model.TeamUser) error {
	return db.Model(&model.TeamUser{}).
		Where("user_id = ? and team_id = ?", teamUser.UserId, teamUser.TeamId).
		Update("role_id", teamUser.RoleId).Error;
}

func (teamService *TeamServiceImpl) MyTeam(userId *int64) (*[]model.Team, error) {
	teams := make([]model.Team, 0, 5)
	err := DB().Where("user_id = ?", *userId).Find(&teams).Error
	return &teams, err
}

func (teamService *TeamServiceImpl) MyJoinTeam(userId *int64) (*[]model.Team, error) {
	teams := make([]model.Team, 0, 5)
	err := DB().Select("team.id, team.name, team.remark, team.user_id").
		Joins("inner join role_user on team.id = role_user.team_id and team.user_id != role_user.user_id").
		Where("team.user_id = ?", *userId).Find(&teams).Error
	return &teams, err
}

func (teamService *TeamServiceImpl) AddTeamUser(teamUser *model.TeamUser) (err error) {
	projectUserMappings := make([]model.ProjectUserMapping, 0, 5)
	if err = DB().Raw("SELECT DISTINCT id, project_id, team_id FROM project_user_mapping WHERE team_id = ?").Scan(&projectUserMappings).Error; err != nil {
		return
	}
	tx := DB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Create(&model.TeamUser{UserId: teamUser.UserId, TeamId: teamUser.TeamId, RoleId: model.ROLE_MEMBER.ID}).Error; err != nil {
		return
	}
	if len(projectUserMappings) == 0 {
		return nil
	}
	proUserDB := make([]model.ProjectUserMapping, 0, len(projectUserMappings))
	for _, val := range projectUserMappings {
		proUserDB = append(proUserDB, model.ProjectUserMapping{TeamId: val.TeamId, UserId:teamUser.UserId, ProjectId: val.ProjectId, AccessLevel: model.ROLE_MEMBER.ID})
	}
	if err = tx.Create(proUserDB).Error; err != nil {
		return err
	}
	return nil
}

func (teamService *TeamServiceImpl) RemoveTeamUser(teamUser *model.TeamUser) (err error) {
	dbTeam := &model.Team{ID: teamUser.TeamId}
	if err = DB().First(dbTeam).Error; err != nil {
		return
	}
	if dbTeam.UserId == teamUser.UserId {
		return errors.New("Team own can not remove self")
	}
	tx := DB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Where("user_id = ? and team_id = ?", teamUser.UserId, teamUser.TeamId).Delete(&model.TeamUser{}).Error; err != nil {
		return 
	}
	if err = tx.Where("user_id = ? and team_id = ?", teamUser.UserId, teamUser.TeamId).Delete(&model.ProjectUserMapping{}).Error; err != nil {
		return
	}
	return nil
}


