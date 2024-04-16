package logic

import (
	"strconv"
	"backbend/models"
	"backbend/dao/redis"
	"go.uber.org/zap"
	
)

/**
 * @Author huchao
 * @Description //TODO 投票功能
 * @Date 11:35 2022/2/14
 **/
 func VoteForPost(userId uint64, p *models.VoteDataForm) error {
	zap.L().Debug("VoteForPost",
		zap.Uint64("userId",userId),
		zap.String("postId", p.PostID),
		zap.Int8("Direction",p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userId)), p.PostID, float64(p.Direction))
}