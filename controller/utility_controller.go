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
	"encoding/json"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap"
	"fmt"
	"beebe/log"
)

var BusMap map[int64]string
var NoLoginResult []byte
var AlreadyLoginResult []byte
var SystemErrorResult []byte

type CommonController struct{}

type UploadResponse struct {
	Url string            `json:"url"`
}

func init() {
	NoLoginResult, _ = json.Marshal(model.ConvertRestResult(model.USER_NO_LOGIN))
	AlreadyLoginResult, _ = json.Marshal(model.ConvertRestResult(model.USER_ALREADY_LOGIN))
	SystemErrorResult, _ = json.Marshal(model.ConvertRestResult(model.SYSTEM_ERROR))
	commonController := new(CommonController)
	Macaron().Group("/common", func() {
		Macaron().Post("/upload", binding.MultipartForm(UploadForm{}), commonController.upload)
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
	if uploadForm.Format == "" {
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


func needLogin(ctx *macaron.Context, sess session.Store) {
	if user := getCurrentUser(sess); user == nil {
		ctx.Resp.Write(NoLoginResult)
	}
}

func noNeedLogin(ctx *macaron.Context, sess session.Store) {
	if user := getCurrentUser(sess); user != nil {
		ctx.Resp.Write(AlreadyLoginResult)
	}
}

func setResponse(ctx *macaron.Context, result interface{}, errCode *model.ErrorCode, err error) {
	if errCode != nil && errCode.Code != model.SUCCESS.Code  {
		size := 3
		if err != nil {
			size = 4
		}
		fields := make([]zapcore.Field, size)
		fields[0] = zap.String("requestUrl", ctx.Req.URL.String())
		fields[1] = zap.Int("errCode", errCode.Code)
		fields[2] = zap.String("errMsg", errCode.Msg)
		if err != nil {
			fields[3] = zap.Error(err)
		}
		str := fmt.Sprintf("request url[%s] log, errMsg: %s;", ctx.Req.URL.String(), errCode.Msg)
		log.Log.Error(str, fields...)
	}
	restResult := model.ConvertRestResult(errCode)
	if errCode.Code == model.SUCCESS.Code && result != nil {
		restResult.SetData(result)
	}
	resultError, erro := json.Marshal(*restResult)
	if erro != nil {
		ctx.Resp.Write(SystemErrorResult)
	}
	ctx.Resp.Write(resultError)
}

func setSuccessResponse(ctx *macaron.Context, result interface{}) {
	setResponse(ctx, result, model.SUCCESS, nil)
}

func setFailResponse(ctx *macaron.Context, errCode *model.ErrorCode, err error) {
	setResponse(ctx, nil, errCode, err)
}

func setErrorResponse(ctx *macaron.Context, errCode *model.ErrorCode)  {
	setResponse(ctx, nil, errCode, nil)
}
