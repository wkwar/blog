package scores

import (
	"math"
	
)

//加速度，决定排名随时间下降的速度快慢
var acceleration float64 = 1.8

//计算帖子的分数 
func Ranking(voteNum, timeDur int64) float64 {
	voteNum -= 1
	return float64(voteNum) / math.Pow(float64(timeDur + 2), acceleration)
}