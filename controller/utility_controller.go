package controller

import (
	"beebe/config"
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/session"
	"beebe/model"
	"beebe/utils"
	"github.com/go-macaron/binding"
	"io"
	"os"
)

var BusMap map[int64]string

type CommonController struct {}

type UploadResponse struct {
	Url 			string			`json:"url"`
}

func init() {
	commonController := new(CommonController)
	Macaron().Group("/common", func() {
		Macaron().Post("/upload",  binding.MultipartForm(UploadForm{}), commonController.upload)
	}, needLogin)
	BusMap = make(map[int64]string)
	BusMap[0] = config.GetConfig().Upload.Default
	BusMap[1] = config.GetConfig().Upload.UserPath
	BusMap[2] = config.GetConfig().Upload.ProjectPath
	BusMap[3] = config.GetConfig().Upload.TeamPath
}

func (commonController *CommonController) upload(uploadForm UploadForm, ctx *macaron.Context, sess session.Store) {
	var busPath string
	if val, ok := BusMap[uploadForm.Bus]; !ok {
		busPath = config.GetConfig().Upload.Default
	} else {
		busPath = val
	}
	if uploadForm.Format== "" {
		uploadForm.Format = "png"
	}
	file, err := uploadForm.ImageUpload.Open()
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	defer file.Close()
	filePath := busPath + utils.GetGuid() + "." + uploadForm.Format
	out, err := os.Create(config.GetConfig().Upload.Base + filePath)
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		setFailResponse(ctx, model.SYSTEM_ERROR, err)
		return
	}
	setSuccessResponse(ctx, &UploadResponse{Url: filePath})
}
