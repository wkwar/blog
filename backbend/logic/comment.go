package logic

import (
	"backbend/models"
	"backbend/dao/mysql"
	"backbend/pkg/snowflake"
	"backbend/dao/redis"
	"go.uber.org/zap"
)

func CreateComment(comment *models.Comment) (err error) {
	// 生成评论ID
	commentID, err := snowflake.GetID()
	if err != nil {
		zap.L().Error("snowflake.GetID() failed", zap.Error(err))
		return
	}

	comment.CommentID = commentID
	//数据库创建
	if err = mysql.CreateComment(comment); err != nil {
		zap.L().Error("mysql.CreateComment(&comment) failed", zap.Error(err))
		return
	}

	//创建成功后，将帖子 的评论数 + 1
	if err = mysql.IncrCommentNum(comment.PostID); err != nil {
		zap.L().Error("mysql.IncrCommentNum(id) failed", zap.Error(err))
		return
	}

	//redis存储  
	if err = redis.CreateComment(comment); err != nil {
		zap.L().Error("redis.CreateComment(comment) failed", zap.Error(err))
		return
	}

	
	


	return
}

func GetCommentListByID(id uint64) (data []*models.Comment, err error) {
	return
}

