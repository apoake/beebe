package service

import (
	"beebe/model"
	"github.com/pkg/errors"
	"github.com/jinzhu/gorm"
)

var projectActionService = new(ProjectActionServiceImpl)
var projectService = new(ProjectServiceImpl)
var workSpaceService = new(WorkSpaceServiceImpl)

func GetProjectActionService() *ProjectActionServiceImpl {
	return projectActionService
}

func GetProjectService() *ProjectServiceImpl {
	return projectService
}

func GetWorkSpaceService() *WorkSpaceServiceImpl {
	return workSpaceService
}

type ProjectService interface{
	AddProject(project *model.Project) error
	AddTeam(projectId *int64, teamId *int64) error
	RemoveTeam(projectId *int64, teamId *int64) error
	UpdateProject(project *model.Project) error
	DeleteProject(project *model.Project) error
	GetProject(projectId *int64) (*model.Project, bool)
	GetProjectsPage(key string, start *int64, limit *int64) (*[]model.Project, error)
	GetProjectByUser(userIds *[]int64) (*[]model.Project, error)
	GetJoiningProjects(userId *int64) (*[]model.Project, error)
}

type ProjectServiceImpl struct{}

func (projectService *ProjectServiceImpl) AddProject(project *model.Project) (err error) {
	tx := DB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Create(project).Error; err != nil {
		return err
	}
	if err = tx.Create(&model.ProjectUserMapping{ProjectId: project.ID, UserId: project.UserId, AccessLevel: model.ROLE_MASTER.ID}).Error; err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (projectService *ProjectServiceImpl) AddTeam(projectId *int64, teamId *int64) (err error) {
	teamUsers := make([]model.TeamUser, 5)
	if !DB().Where("team_id = ?", *teamId).Find(&teamUsers).RecordNotFound() {
		return nil
	}
	tx := DB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	for _, val := range teamUsers {
		if err = tx.Create(&model.ProjectUserMapping{UserId: val.UserId, TeamId: *teamId, ProjectId: *projectId, AccessLevel: val.RoleId}).Error; err != nil {
			return
		}
	}
	tx.Commit()
	return nil
}

func (projectService *ProjectServiceImpl) RemoveTeam(projectId *int64, teamId *int64) error {
	return DB().Where("project_id = ? and team_id = ?", *projectId, *teamId).Delete(model.ProjectUserMapping{}).Error
}

func (projectService *ProjectServiceImpl) UpdateProject(project *model.Project) error {
	dbProject := &model.Project{ID: project.ID}
	DB().First(dbProject)
	dbProject.Name = project.Name
	dbProject.Introduction = project.Introduction
	dbProject.IsPublic = project.IsPublic
	dbProject.ImgUrl = project.ImgUrl
	return DB().Save(dbProject).Error
}

func (projectService *ProjectServiceImpl) DeleteProject(project *model.Project) (err error) {
	dbProject := &model.Project{ID: project.ID}
	DB().First(dbProject)
	if dbProject.UserId != project.UserId {
		return errors.New("only project ownner can delete the project")
	}
	tx := DB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	if err = tx.Delete(dbProject).Error; err != nil {
		return
	}
	projectActions, erro :=GetProjectActionService().GetAllByProjectId(&project.ID)
	if erro != nil {
		return erro
	}
	for _, val := range *projectActions {
		if err = GetProjectActionService().Delete(&val.ActionId, tx); err != nil {
			return
		}
	}
	tx.Commit()
	return nil
}

func (projectService *ProjectServiceImpl) GetProject(projectId int64) (*model.Project, bool) {
	project := new(model.Project)
	isExit := !DB().First(project, projectId).RecordNotFound()
	return project, isExit
}

func (projectService *ProjectServiceImpl) GetProjectsPage(key string, start *int64, limit *int64) (*[]model.Project, error) {
	projects := make([]model.Project, 0, *limit)
	db := DB().Select("name, introduction")
	if key != "" {
		db = db.Where("is_public = 1 and name LIKE ?", "%" + key + "%")
	}
	err := db.Offset(*start).Limit(*limit).Find(&projects).Error
	return &projects, err
}

func (projectService *ProjectServiceImpl) GetProjectByUser(userIds *[]int64) (*[]model.Project, error) {
	projects := make([]model.Project, 0, 5)
	err := DB().Where("user_id in (?)", *userIds).Find(&projects).Error
	return &projects, err
}

func (projectService *ProjectServiceImpl) GetJoiningProjects(userId *int64) (*[]model.Project, error) {
	projects := make([]model.Project, 0, 5)
	err := DB().Select("project.id, project.name, project.introduction").
		Joins("inner join project_user_mapping on project_user_mapping.project_id = project.id").
		Where("project_user_mapping.user_id = ? and project_user_mapping.team_id != 0", userId).Find(&projects).Error
	return &projects, err
}


type ProjectActionService interface {
	Get(actionId *int64) (*model.ProjectAction, bool)
	GetAllByProjectPage(projectId *int64, start *int64, limit *int64) (*[]model.ProjectAction, error)
	GetAllByProjectId(projectId *int64) (*[]model.ProjectAction, error)
	GetByProjectIdAndUrl(projectId *int64, url *string) (*model.ProjectAction, bool)
	CreateProjectAction(projectAction *model.ProjectAction) error
	UpdateProjectAction(projectAction *model.ProjectAction) error
	Delete(actionId *int64) error
}

type ProjectActionServiceImpl struct {}

func (projectActionService *ProjectActionServiceImpl) Get(actionId *int64) (*model.ProjectAction, bool) {
	projectAction := new(model.ProjectAction)
	projectAction.ActionId = *actionId
	isExist := !DB().First(projectAction).RecordNotFound()
	return projectAction, isExist
}

func (projectActionService *ProjectActionServiceImpl) GetAllByProjectPage(projectId *int64, start *int64, limit *int64) (*[]model.ProjectAction, error) {
	projectActions := make([]model.ProjectAction, *limit)
	err := DB().Offset(start).Limit(limit).Where("project_id = ?", projectId).Find(projectActions).Error
	return &projectActions, err
}

func (projectActionService *ProjectActionServiceImpl) GetAllByProjectId(projectId *int64) (*[]model.ProjectAction, error) {
	projectActions := make([]model.ProjectAction, 5)
	err := DB().Where("project_id = ?", *projectId).Find(projectActions).Error
	return &projectActions, err
}

func (projectActionService *ProjectActionServiceImpl) GetByProjectIdAndUrl(projectId *int64, url *string) (*model.ProjectAction, bool) {
	projectAction := &model.ProjectAction{}
	isExist := !DB().Where("project_id = ? and request_url = ?", *projectId, *url).First(projectAction).RecordNotFound()
	return projectAction, isExist
}

func (projectActionService *ProjectActionServiceImpl) CreateProjectAction(projectAction *model.ProjectAction) error {
	return DB().Create(projectAction).Error
}

func (projectActionService *ProjectActionServiceImpl) UpdateProjectAction(projectAction *model.ProjectAction) error {
	dbProjectAction := new(model.ProjectAction)
	dbProjectAction.ActionId = projectAction.ActionId
	if err := DB().Find(dbProjectAction).Find(dbProjectAction).Error; err != nil {
		return err
	}
	dbProjectAction.ActionName = projectAction.ActionName
	dbProjectAction.ActionDesc = projectAction.ActionDesc
	dbProjectAction.RequestType = projectAction.RequestType
	dbProjectAction.RequestUrl = projectAction.RequestUrl
	return DB().Save(dbProjectAction).Error
}

func (projectActionService *ProjectActionServiceImpl) Delete(actionId *int64, db *gorm.DB) (err error) {
	var tx *gorm.DB
	if db == nil {
		tx = DB().Begin()
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()
	} else {
		tx = db
	}
	projectAction := new(model.ProjectAction)
	projectAction.ActionId = *actionId
	tmp := tx.Delete(projectAction)
	if err = tmp.Error; err != nil {
		return err
	}
	if err = paramActionService.DeleteByActionId(actionId, tx); err != nil {
		return err
	}
	return nil
}

type WorkSpaceService interface {
	GetProject(userId *int64) (*[]model.Project, error)
	GetByUserIdAndProjectId(userId *int64, projectId *int64) bool
	Get(userId *int64) (*[]model.WorkSpace, bool)
	AddProject(workSpace *model.WorkSpace) error
	DeleteProject(workSpace *model.WorkSpace) error
}

type WorkSpaceServiceImpl struct {}

func (workspaceService *WorkSpaceServiceImpl) GetProject(userId *int64) (*[]model.Project, error) {
	projects := make([]model.Project, 0, 5)
	err := DB().Select("project.id, project.name, project.img_url, project.introduction, project.is_public").
		Joins("inner join workspace on workspace.project_id = project.id").
		Where("workspace.user_id = ?", *userId).Find(&projects).Error
	return &projects, err
}

func (workspaceService *WorkSpaceServiceImpl) GetByUserIdAndProjectId(userId *int64, projectId *int64) bool {
	return DB().Where("user_id = ? and project_id = ?", *userId, *projectId).RecordNotFound();
}

func (workspaceService *WorkSpaceServiceImpl) Get(userId *int64) (*[]model.WorkSpace, bool) {
	spaces := make([]model.WorkSpace, 0, 5)
	isExist := !db.Where("user_id = ?", *userId).Find(&spaces).RecordNotFound()
	return &spaces, isExist
}

func (workSpaceService *WorkSpaceServiceImpl) AddProject(workSpace *model.WorkSpace) error {
	return DB().Save(workSpace).Error
}

func (workSpaceService *WorkSpaceServiceImpl) DeleteProject(workSpace *model.WorkSpace) error {
	return DB().Where("user_id = ? and project_id = ?", workSpace.UserId, workSpace.ProjectId).Delete(model.WorkSpace{}).Error
}