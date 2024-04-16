package controller

import (
	"backbend/dao/mysql"
	"backbend/logic"
	"backbend/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CommentHandler 创建评论
func CreateReplyHandler(c *gin.Context) {
	var reply models.Reply
	if err := c.BindJSON(&reply); err != nil {
		fmt.Println(err)
		ResponseError(c, CodeInvalidParams)
		return
	}
	fmt.Println("回复 内容", reply)
	
	// 获取作者ID，当前请求的UserID
	user, err := getCurrentUser(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}
	reply.AuthorID = user.UserID
	reply.AuthorName = user.UserName
	reply.CreateTime = time.Now()

	// 创建回复
	err = logic.CreateReply(&reply)
	if(err != nil) {
		zap.L().Error("CreateComment() failed", zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}

	ResponseSuccess(c, nil)
}


func ReplyListHandler(c *gin.Context) {
	idstr := c.Query("comment_id")
	id, _:= strconv.ParseInt(idstr, 10, 64)
	fmt.Println("id", id)
	replys, err := mysql.GetReplyList(id)
	if err != nil {
		zap.L().Error("GetReplyList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, replys)
}
