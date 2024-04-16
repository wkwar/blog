package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/**
 * @Author wkwar
 * @Description //TODO 自定义响应封装,前端根据返回的res.Code来执行不同功能
 * @Date 14:00 2022/3/17
 **/
type ResponseData struct {
	Code    MyCode      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // omitempty当data为空时,不展示这个字段

}

type ResponsePosts struct {
	Code         MyCode      `json:"code"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data,omitempty"`
	Record		interface{} `json:"record"`
	Total        int64       `json:"total"`
	Size         int64       `json:"size"`
	Current_page int64       `json:"current_page"`
	Pages        int64       `json:"pages"` //总页数，
}

func ResponseError(ctx *gin.Context, c MyCode) {
	rd := &ResponseData{
		Code:    c,
		Message: c.Msg(),
		Data:    nil,
	}
	ctx.JSON(http.StatusOK, rd)
}

func ResponseErrorWithMsg(ctx *gin.Context, code MyCode, data interface{}) {
	rd := &ResponseData{
		Code:    code,
		Message: code.Msg(),
		Data:    data,
	}
	ctx.JSON(http.StatusOK, rd)
}

func ResponseSuccess(ctx *gin.Context, data ...interface{}) {
	rd := &ResponseData{
		Code:    CodeSuccess,
		Message: CodeSuccess.Msg(),
		Data:    data,
	}
	ctx.JSON(http.StatusOK, rd)
}

func PostResponseSuccess(ctx *gin.Context, data interface{}, total, current_page, size int64) {
	rd := &ResponsePosts{
		Code:         CodeSuccess,
		Message:      CodeSuccess.Msg(),
		Data:         data,
		Total:        total,
		Current_page: current_page,
		Size:         size,
	}

	ctx.JSON(http.StatusOK, rd)
}
