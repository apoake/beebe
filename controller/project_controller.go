package controller

import (
	"beebe/service"
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/session"
	"beebe/model"
	"github.com/go-macaron/binding"
)


type ProjectController struct {}

func init() {
	projectController := new(ProjectController)
	Macaron().Get("/search", projectController.search)
	Macaron().Group("/project", func() {
		Macaron().Post("/create", binding.Bind(model.Project{}), projectController.create)
		Macaron().Post("/mine", projectController.myProjects)
		Macaron().Post("/join", projectController.myJoiningProjects)
		Macaron().Post("/space", projectController.myWorkspace)
	}, needLogin)
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

func (projectController *ProjectController) create(project model.Project, ctx *macaron.Context, sess session.Store) {
	if project.Name == "" || project.IsPublic == 0 {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	projectdb := &model.Project{IsPublic: project.IsPublic, Introduction:project.Introduction, Name: project.Name, UserId: *getCurrentUserId(sess)}
	err := service.GetProjectService().AddProject(projectdb)
	if err != nil {
		setFailResponse(ctx, model.PROJECT_CREATE_ERROR, err)
		return
	}
	setSuccessResponse(ctx, *projectdb)
}

/**
	myProjects
 */
func (projectController *ProjectController) myProjects(ctx *macaron.Context, sess session.Store) {
	user := getCurrentUser(sess)
	userIds := []int64{user.ID}
	projects, err := service.GetProjectService().GetProjectByUser(&userIds)
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

/**
	myWorkspace
 */
func (projectController *ProjectController) myWorkspace(ctx *macaron.Context, sess session.Store) {
	userId := getCurrentUserId(sess)
	if userId == nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, nil)
		return
	}
	projects, err := service.GetWorkSpaceService().GetProject(userId)
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, projects)
}

