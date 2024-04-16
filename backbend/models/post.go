package models

import (
	"time"
)

/**
 * @Author wkwar
 * @Description //TODO 发帖请求参数
 * @Date 14:00 2023/1/1
 **/
type Post struct {
	PostID      uint64    `json:"post_id,string" db:"post_id"`
	AuthorId    uint64    `json:"author_id" db:"author_id"`
	AuthorName  string    `json:"author_name" db:"author_name"`
	CommunityID uint64    `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	Score		int64	  `json:"score" db:"score"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`

	ViewNum    int64 `json:"view_num" db:"view_num"`       // 查看数量
	CommentNum int64 `json:"comment_num" db:"comment_num"` // 评论数量
	LikeNum    int64 `json:"like_num" db:"like_num"`	   //点赞数量
	UnLikeNum  int64 `json:"unlike_num" db:"unlike_num"`   //不喜欢数量
}

type Search struct {
	KeyWord    string `json:"keyword"`
	TimeRange  int    `json:"timeRange"`
	AuthorName string `json:"author_name"`
}

/**
 * @Author wkwar
 * @Description //TODO 显示贴子请求参数
 * @Date 14:00 2023/1/1
 **/
type ApiPostDetail struct {
	*Post                               // 嵌入帖子结构体
	*CommunityDetail `json:"community"` // 嵌入社区信息
	AuthorName       string             `json:"author_name"`
	VoteNum          int64              `json:"vote_num"`
	Page             int64
	Size             int64
}

type UserLike struct {
	//userID + postID 构成一个
	UserID     uint64    `json:"user_id" db:"user_id"`
	PostID     uint64    `json:"post_id" db:"post_id"`
	Status     int8      `json:"status" db:"status"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
}
