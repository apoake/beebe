package controller

import (
	"github.com/go-macaron/binding"
	"mime/multipart"
	"beebe/config"
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/session"
	"beebe/model"
	"beebe/utils"
	"io"
	"io/ioutil"
)

var BusMap map[int64]string

func init() {
	commonController := new(CommonController)
	Macaron().Group("/common", func() {
		Macaron().Post("/upload", binding.MultipartForm(UploadForm{}), commonController.upload)
	}, needLogin)

	BusMap = make(map[int64]string)
	BusMap[0] = config.GetConfig().Upload.Default
	BusMap[1] = config.GetConfig().Upload.UserPath
	BusMap[2] = config.GetConfig().Upload.ProjectPath
}

type CommonController struct {}

type UploadForm struct {
	bus 			int64					`form:"bus"`
	format        	string                	`form:"format"`
	ImageUpload 	*multipart.FileHeader 	`form:"image"`
}

type UploadResponse struct {
	Url 			string			`json:"url"`
}

func (commonController *CommonController) upload(uploadForm UploadForm, ctx *macaron.Context, sess session.Store) {
	var busPath string
	if val, ok := BusMap[uploadForm.bus]; !ok {
		busPath = config.GetConfig().Upload.Default
	} else {
		busPath = val
	}
	if uploadForm.format == "" {
		uploadForm.format = "png"
	}
	file, err := uploadForm.ImageUpload.Open()
	if err != nil {
		setErrorResponse(ctx, model.SYSTEM_ERROR)
		return
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	filePath := busPath + utils.GetGuid() + "." + uploadForm.format
	if err = ioutil.WriteFile(filePath, b, 0644); err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, &UploadResponse{Url: filePath})
}




