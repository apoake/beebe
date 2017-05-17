package controller

import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/binding"
	"beebe/model"
	"encoding/json"
	"fmt"
	"mime/multipart"
)
type Base struct {}

func (base Base) Error(ctx *macaron.Context, errs binding.Errors) {
	if len(errs) > 0 {
		err := errs[0]
		result := model.ConvertRestResult(model.PARAMETER_INVALID)
		result.Message = fmt.Sprintf("%s: %v", err.Classification, err.FieldNames)
		errBody, _ := json.Marshal(result)
		ctx.Resp.Write(errBody)
	}
}

// userController dto
type UserLogin struct {
	Base
	Account 	string 			`json:"userName" binding:"Required"`
	Password 	string          `json:"password" binding:"Required"`
}

type UserRegister struct {
	Base
	Account  	string			`json:"userName" binding:"Required"`
	Password 	string          `json:"password" binding:"Required"`
	Email    	string          `json:"email" binding:"OmitEmpty;Email"`
}

type UserSearch struct {
	Base
	Account  	string			`json:"userName" binding:"Required"`
}

type UserUpdate struct {
	Base
	Email    	string          `json:"email" binding:"OmitEmpty;Email"`
	ImgUrl  	string			`json:"imgUrl"" binding:"OmitEmpty"`
	Name 		string			`json:"nickName" binding:"OmitEmpty;AlphaDash`
}

type UserPassword struct {
	Base
	Opassword 	string			`json:"opassword" binding:"Required;AlphaDash"`
	Password 	string          `json:"password" binding:"Required;AlphaDash"`
}

type TeamAdd struct {
	Base
	Name 		string		`json:"name" binding:"Required"`
	Remark		string		`json:"remark"`
	LogoUrl		string		`json:"logoUrl"`
}

type Id struct {
	Base
	ID			int64		`json:"id" binding:"Required"`
}

type TeamUpdate struct {
	TeamAdd
	Id
}

type TeamUserDto struct {
	Base
	TeamId			int64		`json:"teamId" binding:"Required"`
	UserId 			int64		`json:"userId" binding:"Required"`
	RoleId			int64		`json:"roleId" binding:"OmitEmpty;Range(1,3)"`
}

// utilityController
type UploadForm struct {
	Base
	Bus 			int64					`form:"bus" binding:"Required;Range(0,3)"`
	Format        	string                	`form:"format" binding:"Required;AlphaDash"`
	ImageUpload 	*multipart.FileHeader 	`form:"image" binding:"Required"`
}

// projectController
type ProjectCreate struct {
	Base
	Name         string        	`json:"name" binding:"Required;AlphaDash"`
	ImgUrl		 string        	`json:"imgUrl"`
	Introduction string         `json:"introduction" binding:"OmitEmpty;AlphaDash"`
	IsPublic     int            `json:"isPublic" binding:"Required;Range(1,2)"`
}

type ProjectUpdate struct {
	ProjectCreate
	Id
}

type ProjectID struct {
	Base
	ProjectId 		int64 		`json:"projectId" binding:"Required"`
}

type ActionID struct {
	Base
	ActionId 		int64		`json:"actionId" binding:"Required"`
}

type ProjectActionCreate struct {
	Base
	ActionName			string		`json:"actionName" binding:"Required;AlphaDash"`
	ActionDesc			string		`json:"actionDesc" binding:"OmitEmpty;AlphaDash"`
	ProjectId			int64		`json:"projectId" binding:"Required"`
	RequestType 		string		`json:"requestType" binding:"Required"`
	RequestUrl			string		`json:"requestUrl" binding:"Required"`
}

type ProjectActionUpdate struct {
	ActionID
	ActionName			string		`json:"actionName" binding:"Required;AlphaDash"`
	ActionDesc			string		`json:"actionDesc" binding:"OmitEmpty;AlphaDash"`
	RequestType 		string		`json:"requestType" binding:"Required"`
	RequestUrl			string		`json:"requestUrl" binding:"Required"`
}