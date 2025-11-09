package controller

import (
	"demo1/internal/app/mydemo/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CODE_SUCCESS         int = 0
	CODE_SYS_ERROR       int = 1
	CODE_PARAMS_ERROR    int = 1001
	CODE_USER_NAME_EXIST int = 2001
	CODE_USER_BLOCKED    int = 2002
	CODE_CHECKIN_REPEAT  int = 3001
)

var CodeMsgMap map[int]string = map[int]string{
	CODE_SUCCESS:         "成功",
	CODE_SYS_ERROR:       "服务出错，稍后再试",
	CODE_PARAMS_ERROR:    "参数错误",
	CODE_USER_NAME_EXIST: "用户名已存在",
	CODE_USER_BLOCKED:    "用户状态异常",
	CODE_CHECKIN_REPEAT:  "重复打卡",
}

func GetMsgByCode(code int) (msg string) {
	msg, ok := CodeMsgMap[code]
	if ok {
		return
	} else {
		msg = CodeMsgMap[CODE_SYS_ERROR]
		return
	}
}

func MakeApiResponse(c *gin.Context, code int, data interface{}) {
	msg := GetMsgByCode(code)
	c.JSON(http.StatusOK, model.APIResponse{
		Code:    code,
		Data:    data,
		Message: msg,
	})
	return
}
