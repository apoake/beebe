package controller

import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/binding"
	"beebe/model"
	"encoding/json"
	"fmt"
)
// userController dto
type UserLogin struct {
	Account 	string 			`json:"userName" binding:"Required"`
	Password 	string          `json:"password" binding:"Required"`
}

func (userLogin UserLogin) Error(ctx *macaron.Context, errs binding.Errors) {
	operateBindError(ctx, errs)
}

type UserRegister struct {
	Account  	string			`json:"userName" binding:"Required"`
	Password 	string          `json:"password" binding:"Required"`
	Email    	string          `json:"email" binding:"OmitEmpty;Email"`
}

func (userRegister UserRegister) Error(ctx *macaron.Context, errs binding.Errors) {
	operateBindError(ctx, errs)
}

type UserSearch struct {
	Account  	string			`json:"userName" binding:"Required"`
}

func (userSearch UserSearch) Error(ctx *macaron.Context, errs binding.Errors) {
	operateBindError(ctx, errs)
}

type UserUpdate struct {
	Email    	string          `json:"email" binding:"OmitEmpty;Email"`
	ImgUrl  	string			`json:"imgUrl binding:"OmitEmpty"`
	Name 		string			`json:"nickName" binding:"OmitEmpty;AlphaDash`
}

func (userUpdate UserUpdate) Error(ctx *macaron.Context, errs binding.Errors) {
	operateBindError(ctx, errs)
}

type UserPassword struct {
	Opassword 	string			`json:"opassword" binding:"Required;AlphaDash"`
	Password 	string          `json:"password" binding:"Required;AlphaDash"`
}

func (userPassword UserPassword) Error(ctx *macaron.Context, errs binding.Errors) {
	operateBindError(ctx, errs)
}

type TeamAdd struct {
	Name 		string		`json:"name" binding:"Required"`
	Remark		string		`json:"remark"`
	LogoUrl		string		`json:"logoUrl"`
}

func (teamAdd TeamAdd) Error(ctx *macaron.Context, errs binding.Errors) {
	operateBindError(ctx, errs)
}

func operateBindError(ctx *macaron.Context, errs binding.Errors) {
	if len(errs) > 0 {
		err := errs[0]
		result := model.ConvertRestResult(model.PARAMETER_INVALID)
		result.Message = fmt.Sprintf("%s: %v", err.Classification, err.FieldNames)
		errBody, _ := json.Marshal(result)
		ctx.Resp.Write(errBody)
	}
}