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
func CreateCommentHandler(c *gin.Context) {
	//获得 帖子id
	postIdStr := c.Query("post_id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param",zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	var comment models.Comment
	if err := c.BindJSON(&comment); err != nil {
		fmt.Println(err)
		ResponseError(c, CodeInvalidParams)
		return
	}
	fmt.Println("评论 内容", comment)

	
	comment.PostID = uint64(postId)
	
	if ok := CheckComment(&comment); ok {
		fmt.Println("评论 违禁词命中")
		ResponseSuccess(c, nil)
		return
	}

	// 获取作者ID，当前请求的UserID
	user, err := getCurrentUser(c)
	if err != nil {
		zap.L().Error("GetCurrentUserID() failed", zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}
	comment.AuthorID = user.UserID
	comment.AuthorName = user.UserName
	comment.CreateTime = time.Now()

	// 创建评论
	err = logic.CreateComment(&comment)
	if(err != nil) {
		zap.L().Error("CreateComment() failed", zap.Error(err))
		ResponseError(c, CodeCreateFailed)
		return
	}

	
	ResponseSuccess(c, nil)
}

// CommentListHandler 评论列表
func CommentListHandler(c *gin.Context) {
	idstr := c.Query("post_id")
	id, _:= strconv.ParseInt(idstr, 10, 64)
	fmt.Println("id", id)
	comments, err := mysql.GetCommentList(id)
	if err != nil {
		zap.L().Error("GetCommentList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, comments)
}

// func CommentListByPostHandler(c *gin.Context) {
// 	ids, ok := c.GetQueryArray("ids")
// 	if !ok {
// 		ResponseError(c, CodeInvalidParams)
// 		return
// 	}
// 	posts, err := mysql.GetCommentListByIDs(ids)
// 	if err != nil {
// 		ResponseError(c, CodeServerBusy)
// 		return
// 	}
// 	ResponseSuccess(c, posts)
// }