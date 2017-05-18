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
		Macaron().Post("/team", binding.Bind(ProjectID{}), projectController.projectTeam)
	}, needLogin)
	Macaron().Group("/space", func() {
		Macaron().Post("/", projectController.myWorkspace)
		Macaron().Post("/addproject", binding.Bind(ProjectID{}), projectController.addWorkspaceProject)
		Macaron().Post("/deleteproject", binding.Bind(ProjectID{}), projectController.deleteWorkspaceProject)
	}, needLogin)
	Macaron().Group("/action", func() {
		Macaron().Post("/create", binding.Bind(ProjectActionCreate{}), projectController.createProjectAction)
		Macaron().Post("/update", binding.Bind(ProjectActionUpdate{}), projectController.updateProjectAction)
		Macaron().Post("/delete", binding.Bind(ActionID{}), projectController.deleteProjectAction)
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
	userId := getCurrentUserId(sess)
	if !service.GetUserService().HasProjectOperateRight(&projectUpdate.ID, &userId) {
		setErrorResponse(ctx, model.USER_NO_RIGHT)
		return
	}
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
func (projectController *ProjectController) deleteProject(id ProjectID, ctx *macaron.Context, sess session.Store) {
	userId := getCurrentUserId(sess)
	if !service.GetUserService().HasProjectOperateRight(&id.ProjectId, &userId) {
		setErrorResponse(ctx, model.USER_NO_RIGHT)
		return
	}
	if err := service.GetProjectService().DeleteProject(&model.Project{ID: id.ProjectId, UserId: userId}); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}


type Wspace struct {
	model.Project
	InWspace 		bool		`json:"isWspace"`
}

func convertProjectToWspace(projects *[]model.Project, userId *int64) *[]Wspace {
	length := len(*projects)
	resultSpace := make([]Wspace, length, length)
	if length == 0 {
		return &resultSpace
	}
	spaces, isExist := service.GetWorkSpaceService().Get(userId)
	isLoop := isExist && len(*spaces) == 0
	for _, val := range *projects {
		isIn := false
		if isLoop {
			for _, w := range *spaces {
				if w.ProjectId == val.ID {
					isIn = true
					break
				}
			}
		}
		resultSpace = append(resultSpace, Wspace{InWspace: isIn, Project: val})
	}
	return &resultSpace
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
	setSuccessResponse(ctx, convertProjectToWspace(projects, &user.ID))
}

/**
	myJoiningProjects
 */
func (projectController *ProjectController) myJoiningProjects(ctx *macaron.Context, sess session.Store) {
	userId := getCurrentUserId(sess)
	projects, err := service.GetProjectService().GetJoiningProjects(&userId)
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, convertProjectToWspace(projects, &userId))
}

/**
	项目中的团队
 */
func (projectController *ProjectController) projectTeam(id ProjectID, ctx *macaron.Context, sess session.Store){
	userId := getCurrentUserId(sess)
	if !service.GetUserService().HasProjectSeeRight(&id.ProjectId, &userId) {
		setErrorResponse(ctx, model.USER_NO_RIGHT)
		return
	}
	teams, err := service.GetTeamService().GetTeamByProjectId(&id.ProjectId)
	if err != nil {
		setErrorResponse(ctx, model.SYSTEM_ERROR)
		return
	}
	setSuccessResponse(ctx, teams)
}

/**
	项目中添加团队
 */
func (projectController *ProjectController) addProjectTeam(projectTeamAdd ProjectTeamAdd, ctx *macaron.Context, sess session.Store) {
	userId := getCurrentUserId(sess)
	if !service.GetUserService().HasProjectOperateRight(&projectTeamAdd.ProjectId, &userId) {
		setErrorResponse(ctx, model.USER_NO_RIGHT)
		return
	}
	if _, ok := service.GetTeamService().Get(&projectTeamAdd.TeamId.TeamId); !ok {
		setErrorResponse(ctx, model.TEAM_NOT_EXIST)
		return
	}
	if err := service.GetProjectService().AddTeam(&projectTeamAdd.ProjectId, &userId); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

func (projectController *ProjectController) removeProjectTeam(projectTeamRemove ProjectTeamRemove, ctx *macaron.Context, sess session.Store) {
	userId := getCurrentUserId(sess)
	if !service.GetUserService().HasProjectOperateRight(&projectTeamRemove.ProjectId, &userId) {
		setErrorResponse(ctx, model.USER_NO_RIGHT)
		return
	}
	if _, ok := service.GetTeamService().Get(&projectTeamRemove.TeamId.TeamId); !ok {
		setErrorResponse(ctx, model.TEAM_NOT_EXIST)
		return
	}
	if err := service.GetProjectService().RemoveTeam(&projectTeamRemove.ProjectId, &userId); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
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
	userId := getCurrentUserId(sess)
	if ok := service.GetWorkSpaceService().GetByUserIdAndProjectId(&userId, &projectId.ProjectId); !ok {
		setErrorResponse(ctx, model.WORKSPACE_ALREADY_IN)
		return
	}
	if err := service.GetWorkSpaceService().AddProject(&model.WorkSpace{UserId: userId, ProjectId: projectId.ProjectId}); err != nil {
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
func (projectController *ProjectController) createProjectAction(projectActionCreate ProjectActionCreate, ctx *macaron.Context, sess session.Store) {
	if err := service.GetProjectActionService().CreateProjectAction(&model.ProjectAction{ProjectId: projectActionCreate.ProjectId, RequestUrl: projectActionCreate.RequestUrl,
		RequestType: projectActionCreate.RequestType, ActionName: projectActionCreate.ActionName, ActionDesc: projectActionCreate.ActionDesc}); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

func (projectController *ProjectController) updateProjectAction(projectActionUpdate ProjectActionUpdate, ctx *macaron.Context, sess session.Store) {
	userId := getCurrentUserId(sess)
	if !service.GetUserService().HasProjectSeeRightByActionId(&projectActionUpdate.ActionId, &userId) {
		setErrorResponse(ctx, model.USER_NO_RIGHT)
		return
	}
	if err := service.GetProjectActionService().UpdateProjectAction(&model.ProjectAction{ActionId: projectActionUpdate.ActionId, ActionName: projectActionUpdate.ActionName,
		ActionDesc: projectActionUpdate.ActionDesc, RequestType: projectActionUpdate.RequestType, RequestUrl: projectActionUpdate.RequestUrl}); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}

func (projectController *ProjectController) deleteProjectAction(actionId ActionID, ctx *macaron.Context, sess session.Store) {
	userId := getCurrentUserId(sess)
	if !service.GetUserService().HasProjectSeeRightByActionId(&actionId.ActionId, &userId) {
		setErrorResponse(ctx, model.USER_NO_RIGHT)
		return
	}
	if err := service.GetProjectActionService().Delete(&actionId.ActionId, nil); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, nil)
}
