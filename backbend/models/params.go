package models

/**
 * @Author wkwar
 * @Description //TODO magic string
 * @Date 21:56 2022/2/15
 **/
 const (
	OrderTime = "time"
	OrderScore = "score"
)

/**
 * @Author wkwar
 * @Description //TODO 获取帖子列表query string参数
 * @Date 20:20 2022/3/18
 **/
 type ParamPostList struct {
	CommunityID uint64  `json:"community_id" form:"community_id"`  // 可以为空
	Page  int64			`json:"page" form:"page"`				   // 页码
	Size  int64			`json:"size" form:"size"`				   // 每页数量
	Order string		`json:"order" form:"order" example:"score"`// 排序依据
}