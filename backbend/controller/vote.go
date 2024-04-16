package controller

import (
	"backbend/models"
	"backbend/logic"
	"go.uber.org/zap"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)



func VoteHandler(c *gin.Context) {
	// 参数校验,给哪个文章投什么票
	vote := new(models.VoteDataForm)
	if err := c.ShouldBindJSON(&vote); err != nil {
		errs, ok := err.(validator.ValidationErrors)	// 类型断言
		if !ok {
			ResponseError(c, CodeInvalidParams)
			return
		}
		errdata := removeTopStruct(errs.Translate(trans))	// 翻译并去除掉错误提示中的结构体标识
		ResponseErrorWithMsg(c, CodeInvalidParams, errdata)
		return
	}
	// 获取当前请求用户的id
	user, err := getCurrentUser(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	// 具体投票的业务逻辑
	userID := user.UserID
	if err := logic.VoteForPost(userID, vote); err != nil {
		zap.L().Error("logic.VoteForPost() failed",zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}