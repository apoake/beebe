package project

import (
	"beebe/model/project"
	"beebe/service"
)

type ProjectService interface{
	AddProject(project *project.Project) error
	UpdateProject(project *project.Project) error
	DeleteProject(projectId *int64) error
	GetProject(projectId *int64) *project.Project
	GetProjectsPage(start *int64, limit *int64) *[]project.Project
}

type ProjectServiceImpl struct{}

func (projectService *ProjectServiceImpl) AddProject(project *project.Project) error {
	tx := service.DB().Begin()
	if err := tx.Create(project).Error; err != nil {
		tx.Rollback()
		return err
	}
	// TODO 处理初始化项目与用户的对应关系
	return tx.Commit().Error
}

func (projectService *ProjectServiceImpl) UpdateProject(project *project.Project) error {
	dbProject := new(project.Project)
	service.DB().First(dbProject)
	dbProject.Name = project.Name
	dbProject.Introduction = project.Introduction
	dbProject.IsPublic = project.IsPublic
	return service.DB().Save(dbProject).Error
}

func (projectService *ProjectServiceImpl) DeleteProject(projectId *int64) error {
	project := new(project.Project)
	project.ID = *projectId
	return service.DB().Delete(project).Error
}

func (projectService *ProjectServiceImpl) GetProject(projectId int64) (*project.Project, error) {
	project := new(project.Project)
	err := service.DB().First(project, projectId).Error
	return project, err
}

func (projectService *ProjectServiceImpl) GetProjectsPage(start *int64, limit *int64) (*[]project.Project, error) {
	projects := make([]project.Project, *limit)
	err := service.DB().Offset(start).Limit(limit).Find(projects).Error
	return projects, err
}


