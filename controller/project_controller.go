package controller

import (
	"beebe/service"
	"github.com/eleme/taco/log"
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/session"
	"beebe/model"
)


type ProjectController struct {}

func init() {
	projectController := new(ProjectController)
	Macaron().Get("/search", projectController.search, jsonResponse)
}

/**
 *	Search project
 */
func (projectController *ProjectController) search(ctx *macaron.Context) {
	search := ctx.Query("search")
	pageNo := ctx.QueryInt64("pageNo")
	limit := ctx.QueryInt64("pageSize")
	var start int64
	if pageNo < 1 {
		start = 0
	} else {
		start = (pageNo-1) * limit
	}
	result, err := service.GetProjectService().GetProjectsPage(search, &start, &limit)
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, result)
}

/**
	myProjects
 */
func (projectController *ProjectController) myProjects(ctx *macaron.Context, sess session.Store) {
	user := getCurrentUser(sess)
	userIds := []int64{user.ID}
	projects, err := service.GetProjectService().GetProjectByUser(userIds)
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, projects)
}

/**
	myJoiningProjects
 */
func (projectController *ProjectController) myJoiningProjects(ctx *macaron.Context, sess session.Store) {
	user := getCurrentUser(sess)
	projects, err := service.GetProjectService().GetJoiningProjects(&user.ID)
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
	}
	setSuccessResponse(ctx, projects)
}


