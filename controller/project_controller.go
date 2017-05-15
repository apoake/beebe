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
		Macaron().Post("/create", binding.Bind(model.Project{}), projectController.createProject)
		Macaron().Post("/update", binding.Bind(model.Project{}), projectController.updateProject)
		Macaron().Post("/delete", binding.Bind(model.Project{}), projectController.deleteProject)
		Macaron().Post("/mine", projectController.myProjects)
		Macaron().Post("/join", projectController.myJoiningProjects)
	}, needLogin)
	Macaron().Group("/space", func() {
		Macaron().Post("/", projectController.myWorkspace)
		Macaron().Post("/addproject", binding.Bind(model.WorkSpace{}), projectController.addWorkspaceProject)
		Macaron().Post("/deleteproject", binding.Bind(model.WorkSpace{}), projectController.deleteWorkspaceProject)
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
func (projectController *ProjectController) createProject(project model.Project, ctx *macaron.Context, sess session.Store) {
	if project.Name == "" || project.IsPublic == 0 {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	projectdb := &model.Project{IsPublic: project.IsPublic, Introduction: project.Introduction, Name: project.Name, UserId: getCurrentUserId(sess)}
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
func (projectController *ProjectController) updateProject(project model.Project, ctx *macaron.Context, sess session.Store) {
	if project.ID == 0 || project.Name == "" || project.IsPublic == 0 {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	err := service.GetProjectService().UpdateProject(&project)
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

/**
	delete project
 */
func (projectController *ProjectController) deleteProject(project model.Project, ctx *macaron.Context, sess session.Store) {
	if project.ID == 0 {
		setErrorResponse(ctx, model.PARAMETER_INVALID)
		return
	}
	project.UserId = getCurrentUserId(sess)
	err := service.GetProjectService().DeleteProject(&project)
	if err != nil {
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
	if userId, ok := getUserId(ctx, sess); ok {
		projects, err := service.GetProjectService().GetJoiningProjects(&userId)
		if err != nil {
			setFailResponse(ctx, model.SYSTEM_ERROR, err)
		}
		setSuccessResponse(ctx, projects)
	}
}

/**
	myWorkspace
 */
func (projectController *ProjectController) myWorkspace(ctx *macaron.Context, sess session.Store) {
	if userId, ok := getUserId(ctx, sess); ok {
		projects, err := service.GetWorkSpaceService().GetProject(&userId)
		if err != nil {
			setFailResponse(ctx, model.SYSTEM_ERROR, err)
			return
		}
		setSuccessResponse(ctx, projects)
	}
}

func (projectController *ProjectController) addWorkspaceProject(workspace model.WorkSpace, ctx *macaron.Context, sess session.Store) {
	if userId, ok := getUserId(ctx, sess); ok {
		if workspace.ProjectId == 0 {
			setErrorResponse(ctx, model.PARAMETER_INVALID)
			return
		}
		err := service.GetWorkSpaceService().AddProject(&model.WorkSpace{UserId: userId, ProjectId: workspace.ProjectId})
		if err != nil {
			setFailResponse(ctx, model.SYSTEM_ERROR, err)
			return
		}
		setSuccessResponse(ctx, nil)
	}
}

func (projectController *ProjectController) deleteWorkspaceProject(workspace model.WorkSpace, ctx *macaron.Context, sess session.Store) {
	if userId, ok := getUserId(ctx, sess); ok {
		if workspace.ProjectId == 0 {
			setErrorResponse(ctx, model.PARAMETER_INVALID)
			return
		}
		err := service.GetWorkSpaceService().DeleteProject(&model.WorkSpace{UserId: userId, ProjectId: workspace.ProjectId})
		if err != nil {
			setFailResponse(ctx, model.SYSTEM_ERROR, err)
			return
		}
		setSuccessResponse(ctx, nil)
	}
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
