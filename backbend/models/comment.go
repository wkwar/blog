package models

import (
	"time"
)

/**
 * @Author wkwar
 * @Description //TODO 创建评论参数
 * @Date 14:00 2023/1/1
 **/
type Comment struct {
	PostID     uint64    `db:"post_id" json:"post_id"`       //归属于哪个帖子
	CommentID  uint64    `db:"comment_id" json:"comment_id"` //评论ID
	AuthorID   uint64    `db:"author_id" json:"author_id"`   //评论作者
	AuthorName string    `db:"author_name" json:"author_name"`
	Content    string    `db:"content" json:"content"`
	CreateTime time.Time `db:"create_time" json:"create_time"`

	ReplyNum uint64   `db:"reply_num" json:"reply_num"` //评论回复数
	LikeNum  uint64   `db:"like_num" json:"like_num"`   //点赞数
	Replys   []*Reply `json:"replys"`
}

//属于哪一个评论下的 回复
type Reply struct {
	CommentID  uint64    `db:"comment_id" json:"comment_id"` //评论ID
	AuthorID   uint64    `db:"author_id" json:"author_id"`
	AuthorName string    `db:"author_name" json:"author_name"`
	ParentID   uint64    `db:"parent_id" json:"parent_id"`
	ParentName string    `db:"parent_name" json:"parent_name"`
	Content    string    `db:"content" json:"content"`
	CreateTime time.Time `db:"create_time" json:"create_time"`
	ReplyNum uint64 ` json:"reply_num"` //评论回复数
	LikeNum  uint64 ` json:"like_num"`  //点赞数

}
