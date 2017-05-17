package controller

import (
	"beebe/service"
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/session"
	"beebe/model"
	"github.com/go-macaron/binding"
)

type ProjectController struct{}

func init() {
	projectController := new(ProjectController)
	Macaron().Get("/search", projectController.search)
	Macaron().Group("/project", func() {
		Macaron().Post("/create", binding.Bind(ProjectCreate{}), projectController.createProject)
		Macaron().Post("/update", binding.Bind(ProjectUpdate{}), projectController.updateProject)
		Macaron().Post("/delete", binding.Bind(Id{}), projectController.deleteProject)
		Macaron().Post("/mine", projectController.myProjects)
		Macaron().Post("/join", projectController.myJoiningProjects)
	}, needLogin)
	Macaron().Group("/space", func() {
		Macaron().Post("/", projectController.myWorkspace)
		Macaron().Post("/addproject", binding.Bind(ProjectID{}), projectController.addWorkspaceProject)
		Macaron().Post("/deleteproject", binding.Bind(ProjectID{}), projectController.deleteWorkspaceProject)
	}, needLogin)
	Macaron().Group("/action", func() {
		Macaron().Post("/create", binding.Bind(model.Project{}), projectController.createProjectAction)
		Macaron().Post("/update", binding.Bind(model.Project{}), projectController.updateProjectAction)
		Macaron().Post("/delete", binding.Bind(model.Project{}), projectController.deleteProjectAction)
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
		start = (pageNo - 1) * limit
	}
	result, err := service.GetProjectService().GetProjectsPage(search, &start, &limit)
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, result)
}

/**
	create project
 */
func (projectController *ProjectController) createProject(projectCreate ProjectCreate, ctx *macaron.Context, sess session.Store) {
	projectdb := &model.Project{IsPublic: projectCreate.IsPublic, Introduction: projectCreate.Introduction, Name: projectCreate.Name, ImgUrl: projectCreate.ImgUrl, UserId: getCurrentUserId(sess)}
	err := service.GetProjectService().AddProject(projectdb)
	if err != nil {
		setFailResponse(ctx, model.PROJECT_CREATE_ERROR, err)
		return
	}
	setSuccessResponse(ctx, *projectdb)
}

/**
	update project
 */
func (projectController *ProjectController) updateProject(projectUpdate ProjectUpdate, ctx *macaron.Context, sess session.Store) {
	if err := service.GetProjectService().UpdateProject(&model.Project{ID: projectUpdate.ID, ImgUrl:projectUpdate.ImgUrl, Name:projectUpdate.Name,
		IsPublic: projectUpdate.IsPublic, Introduction: projectUpdate.Introduction}); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

/**
	delete project
 */
func (projectController *ProjectController) deleteProject(id Id, ctx *macaron.Context, sess session.Store) {
	if err := service.GetProjectService().DeleteProject(&model.Project{ID: id.ID, UserId: getCurrentUserId(sess)}); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
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
	userId := getCurrentUserId(sess)
	projects, err := service.GetProjectService().GetJoiningProjects(&userId)
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
	projects, err := service.GetWorkSpaceService().GetProject(&userId)
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, projects)
}

func (projectController *ProjectController) addWorkspaceProject(projectId ProjectID, ctx *macaron.Context, sess session.Store) {
	if err := service.GetWorkSpaceService().AddProject(&model.WorkSpace{UserId: getCurrentUserId(sess), ProjectId: projectId.ProjectId}); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

func (projectController *ProjectController) deleteWorkspaceProject(projectId ProjectID, ctx *macaron.Context, sess session.Store) {
	if err := service.GetWorkSpaceService().DeleteProject(&model.WorkSpace{UserId: getCurrentUserId(sess), ProjectId: projectId.ProjectId}); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

/**
	project action
 */
func (projectController *ProjectController) createProjectAction(projectAction model.ProjectAction, ctx *macaron.Context, sess session.Store) {
	if projectAction.ProjectId == 0 || projectAction.ActionName == "" ||
		projectAction.RequestType == "" || projectAction.RequestUrl == "" {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	err := service.GetProjectActionService().CreateProjectAction(&projectAction)
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

func (projectController *ProjectController) updateProjectAction(projectAction model.ProjectAction, ctx *macaron.Context, sess session.Store) {
	if projectAction.ProjectId == 0 || projectAction.ActionName == "" ||
		projectAction.RequestType == "" || projectAction.RequestUrl == "" {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	userId := getCurrentUserId(sess)
	if !service.GetUserService().HasProjectRight(&projectAction.ProjectId, &userId) {
		setErrorResponse(ctx, model.USER_NO_RIGHT)
		return
	}
	err := service.GetProjectActionService().UpdateProjectAction(&projectAction)
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

func (projectController *ProjectController) deleteProjectAction(projectAction model.ProjectAction, ctx *macaron.Context, sess session.Store) {
	if projectAction.ActionId == 0 {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	userId := getCurrentUserId(sess)
	if !service.GetUserService().HasProjectRight(&projectAction.ProjectId, &userId) {
		setErrorResponse(ctx, model.USER_NO_RIGHT)
		return
	}
	if err := service.GetProjectActionService().Delete(&projectAction.ActionId); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}
