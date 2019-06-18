package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.yitum.com/saas/shop-admin/model/trans"
)

const (
	MsgOk       = 0   // MsgOk ajax输出错误码，成功
	MsgRedirect = 302 // MsgRedirect ajax输出错误码，重定向
	MsgErr      = 1   // MsgErr 错误

	MsgLoginUserNil           = 10001
	MsgAuthenticateSessionErr = 10002
	MsgLoginReqErr            = 10003

	MsgRegisterReqParamErr      = 10004
	MsgRegisterNnIsNilErr       = 10005
	MsgRegisterNnIsInvalidErr   = 10006
	MsgRegisterPwdIsNilErr      = 10007
	MsgRegisterRePwdIsNilErr    = 10008
	MsgRegisterPwdNotEqRePwdErr = 10009
	MsgRegisterPwdLengthErr     = 10010
	MsgRegisterCreateUserErr    = 10011

	MsgAccountUserInfoErr = 10012
	MsgAppListParamErr    = 10013
	MsgAppAddParamErr     = 10014
	MsgAppAddCreateErr    = 10015
	MsgAppUpdateParamErr  = 10016
	MsgAppUpdateCreateErr = 10017

	MsgUpdateReqParamErr      = 10018
	MsgUpdateNnIsNilErr       = 10019
	MsgUpdateNnIsInvalidErr   = 10020
	MsgUpdatePwdIsNilErr      = 10021
	MsgUpdateRePwdIsNilErr    = 10022
	MsgUpdatePwdNotEqRePwdErr = 10023
	MsgUpdatePwdLengthErr     = 10024
	MsgUpdateCreateUserErr    = 10025

	MsgAccountListParamErr = 10026

	MsgAuthenticateOauth2Err = 10027
)

// JSONOut 统一JSON格式输出。
type JSONOut struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Result  interface{} `json:"result"`
}

// JSON 提供了系统标准JSON输出方法。
func JSON(c *gin.Context, Code int, message string, result ...interface{}) {
	j := new(JSONOut)
	j.Code = Code
	j.Message = message

	if len(result) > 0 {
		j.Result = result[0]
	} else {
		j.Result = ""
	}
	c.JSON(http.StatusOK, j)
	return
}

// JSONOK 返回正确响应，并提供result可选参数
func JSONOK(c *gin.Context, result ...interface{}) {
	j := new(JSONOut)
	j.Code = 0
	j.Message = "ok"
	if len(result) > 0 {
		j.Result = result[0]
	} else {
		j.Result = ""
	}
	c.JSON(http.StatusOK, j)
	return
}

// JSON 提供了输出List的方法。
func JSONList(c *gin.Context, data interface{}, total int) {
	j := new(JSONOut)
	j.Code = 0
	j.Message = "ok"
	j.Result = trans.RespList{
		List:        data,
		TotalNumber: total,
	}
	c.JSON(http.StatusOK, j)
	return
}

// JSONErr 返回错误响应。
func JSONErr(c *gin.Context) {
	j := new(JSONOut)
	j.Code = 1
	j.Message = "err"
	c.JSON(http.StatusOK, j)
	return
}
