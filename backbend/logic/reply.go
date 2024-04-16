package logic

import (
	"backbend/models"
	"backbend/dao/mysql"
	"backbend/dao/redis"
	"go.uber.org/zap"
)

func CreateReply(reply *models.Reply) (err error) {
	//数据库创建
	if err = mysql.CreateReply(reply); err != nil {
		zap.L().Error("mysql.CreateReply(reply) failed", zap.Error(err))
		return
	}

	//创建成功后，将评论 的回复数 + 1
	if err = mysql.IncrReplytNum(reply.CommentID); err != nil {
		zap.L().Error("mysql.IncrReplytNum(id) failed", zap.Error(err))
		return
	}
	
	//redis存储  
	if err = redis.CreateReply(reply); err != nil {
		zap.L().Error("redis.CreateReply(comment) failed", zap.Error(err))
		return
	}
	return
}



