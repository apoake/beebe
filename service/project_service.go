package service

import (
	"beebe/model"
	"github.com/pkg/errors"
)

var projectActionService = new(ProjectActionServiceImpl)
var projectService = new(ProjectServiceImpl)

func ProjectActionService() *ProjectActionServiceImpl {
	return projectActionService
}

func ProjectService() *ProjectServiceImpl {
	return projectService
}

type ProjectService interface{
	AddProject(project *model.Project) error
	UpdateProject(project *model.Project) error
	DeleteProject(projectId *int64) error
	GetProject(projectId *int64) (*model.Project, error)
	GetProjectsPage(start *int64, limit *int64) (*[]model.Project, error)
}

type ProjectServiceImpl struct{}

func (projectService *ProjectServiceImpl) AddProject(project *model.Project) error {
	tx := DB().Begin()
	if err := tx.Create(project).Error; err != nil {
		tx.Rollback()
		return err
	}
	// TODO 处理初始化项目与用户的对应关系
	return tx.Commit().Error
}

func (projectService *ProjectServiceImpl) UpdateProject(project *model.Project) error {
	dbProject := new(model.Project)
	DB().First(dbProject)
	dbProject.Name = project.Name
	dbProject.Introduction = project.Introduction
	dbProject.IsPublic = project.IsPublic
	return DB().Save(dbProject).Error
}

func (projectService *ProjectServiceImpl) DeleteProject(projectId *int64) error {
	project := new(model.Project)
	project.ID = *projectId
	return DB().Delete(project).Error
}

func (projectService *ProjectServiceImpl) GetProject(projectId int64) (*model.Project, error) {
	project := new(model.Project)
	err := DB().First(project, projectId).Error
	return project, err
}

func (projectService *ProjectServiceImpl) GetProjectsPage(start *int64, limit *int64) (*[]model.Project, error) {
	projects := make([]model.Project, *limit)
	err := DB().Offset(start).Limit(limit).Find(projects).Error
	return projects, err
}

type ProjectActionService interface {
	Get(actionId *int64) (*model.ProjectAction, error)
	GetAllByProjectPage(projectId *int64, start *int64, limit *int64) (*[]model.ProjectAction, error)
	CreateProjectAction(projectAction *model.ProjectAction) error
	UpdateProjectAction(projectAction *model.ProjectAction) error
	Delete(actionId *int64) error
}

type ProjectActionServiceImpl struct {}

func (projectActionService *ProjectActionServiceImpl) Get(actionId *int64) (*model.ProjectAction, error) {
	projectAction := new(model.ProjectAction)
	projectAction.ActionId = *actionId
	err := DB().First(projectAction).Error
	return projectAction, err
}

func (projectActionService *ProjectActionServiceImpl) GetAllByProjectPage(projectId *int64, start *int64, limit *int64) (*[]model.ProjectAction, error) {
	projectActions := make([]model.ProjectAction, *limit)
	err := DB().Offset(start).Limit(limit).Where("project_id = ?", projectId).Find(projectActions).Error
	return projectActions, err
}

func (projectActionService *ProjectActionServiceImpl) CreateProjectAction(projectAction *model.ProjectAction) error {
	return DB().Create(projectAction).Error
}

func (projectActionService *ProjectActionServiceImpl) UpdateProjectAction(projectAction *model.ProjectAction) error {
	dbProjectAction := new(model.ProjectAction)
	dbProjectAction.ActionId = projectAction.ActionId
	DB().Find(dbProjectAction).Find()
	dbProjectAction.ActionName = projectAction.ActionName
	dbProjectAction.ActionDesc = projectAction.ActionDesc
	dbProjectAction.RequestType = projectAction.RequestType
	dbProjectAction.RequestUrl = projectAction.RequestUrl
	return DB().Save(dbProjectAction).Error
}

func (projectActionService *ProjectActionServiceImpl) Delete(actionId *int64) (err error) {
	tx := DB().Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	projectAction := new(model.ProjectAction)
	projectAction.ActionId = *actionId
	tmp := tx.Delete(projectAction)
	if err = tmp.Error; err != nil {
		return err
	}
}