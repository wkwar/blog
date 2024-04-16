package controller

import (
	"backbend/logic"
	"backbend/models"
	"backbend/models/constants"
	"strconv"
	"time"
	"go.uber.org/zap"
	"github.com/gin-gonic/gin"
)


func UserLikeHandler(c *gin.Context) {
	var u models.UserLike
	postIdStr := c.Param("id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param",zap.Error(err))
		ResponseError(c,CodeInvalidParams)
	}
	
	u.PostID = uint64(postId)
	//处理 点赞请求
	user, err := getCurrentUser(c)
	if err != nil {
		zap.L().Error("getCurrentUser failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	u.UserID = user.UserID
	u.CreateTime = time.Now()
	u.Status = constants.NotLikedYet
	err = logic.UserLike(&u)
	if err != nil {
		zap.L().Error("logic.UserLike failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}


	// 3、返回响应  
	ResponseSuccess(c, u)
}