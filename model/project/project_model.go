package project

import "beebe/model"

type Project struct {
	model.Model
	ID           int64			`gorm:"primary_key" json:"id"`
	Version      string        	`grom:"column:"version" json:"version"`
	Name         string        	`grom:"column:"name" json:"name"`
	UserId       int64        	`grom:"column:"user_id json:"userId"`
	Introduction string         `grom:"column:"introduction" json:"introduction"`
	IsPublic     int            `grom:"column:"is_public" json:"isPublic"`
	ProjectData  string         `grom:"column:"project_data" json:"projectData"`
	MockNum      int32          `grom:"column:"mock_num" json:"mock_num"`
}

func (Project) TableName() string {
	return "project"
}

type ProjectUserMapping struct {
	model.Model
	ProjectId   int64        `grom:"column:"project_id" json:"projectId"`
	TeamId		int16		 `grom:"column:"team_id" json:"teamId"`
	UserId      int64        `grom:"column:"user_id" json:"userId"`
	AccessLevel int64        `grom:"column:"access_level" json:"accessLevel"`
}

func (ProjectUserMapping) TableName() string {
	return "project_user_mapping"
}