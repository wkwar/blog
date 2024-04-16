package redis

import (
	"math"
	"time"
	"github.com/go-redis/redis/v8"
)

/**
** 通过点赞，评论数量 来 定义排行榜 排行榜实现动态更新
**/

const (
	OneWeekInSeconds         = 7 * 24 * 3600
	VoteScore        float64 = 432	// 每一票的值432分
	PostPerAge               = 20
)


func VoteForPost(userID string, postID string, v float64) (err error) {
	// 去redis取帖子发布时间
	postTime := client.ZScore(ctx, KeyPostTimeZSet, postID).Val()
	//判断帖子时间 --- 如果过期了就移除
	if float64(time.Now().Unix()) - postTime > OneWeekInSeconds {		// Unix()时间戳
		// 不允许投票了
		return ErrorVoteTimeExpire
	}
	// 2、更新帖子的分数
	// 2和3 需要放到一个pipeline事务中操作
	// 判断是否已经投过票 查当前用户给当前帖子的投票记录
	key := KeyPostVotedZSetPrefix + postID
	ov := client.ZScore(ctx, key, userID).Val()

	// 更新：如果这一次投票的值和之前保存的值一致，就提示不允许重复投票
	if v == ov {
		return ErrVoteRepeted
	}
	var op float64
	if v > ov {
		op = 1
	}else {
		op = -1
	}
	diffAbs := math.Abs(ov - v)		// 计算两次投票的差值
	pipeline := client.TxPipeline()	// 事务操作
	_, err = pipeline.ZIncrBy(ctx, KeyPostScoreZSet, VoteScore*diffAbs*op, postID).Result() // 更新分数
	if ErrorVoteTimeExpire != nil {
		return err
	}
	// 3、记录用户为该帖子投票的数据
	if v == 0 {
		_, err = client.ZRem(ctx, key, postID).Result()
	} else {
		pipeline.ZAdd(ctx, key, &redis.Z{ // 记录已投票
			Score:  v,		// 赞成票还是反对票
			Member: userID,
		})
	}

	_, err = pipeline.Exec(ctx)
	return err
} 

