package service

import (
	"beebe/model"
)

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


