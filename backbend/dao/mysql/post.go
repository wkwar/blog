package mysql

import (
	"strings"
	"database/sql"
	"backbend/models"
	"go.uber.org/zap"
	"github.com/jmoiron/sqlx"
)

/**
 * @Author wkwar
 * @Description //TODO 根据Id查询帖子详情
 * @Date 21:53 2022/3/18
 **/
// CreatePost 创建帖子
func CreatePost(post *models.Post) (err error) {
	sqlStr := `insert into post(
	post_id, title, content, author_id, author_name, community_id)
	values(?,?,?,?,?,?)`
	_, err = db.Exec(sqlStr, post.PostID, post.Title,
		post.Content, post.AuthorId, post.AuthorName, post.CommunityID)
	if err != nil {
		zap.L().Error("insert post failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

func UpdatePost(post *models.Post) (err error) {
	sqlStr := `update post set title = ?, content = ?, community_id = ?, view_num = ? where post_id = ?`
	_, err = db.Exec(sqlStr, post.Title,
		post.Content, post.CommunityID, post.ViewNum, post.PostID)
	if err != nil {
		zap.L().Error("insert post failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

func DeletePost(postID uint64) (err error) {
	sqlStr := `delete from post where post_id = ?`
	_, err = db.Exec(sqlStr, postID)
	if err != nil {
		zap.L().Error("insert post failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

/**
 * @Author wkwar
 * @Description //TODO 根据Id查询帖子详情
 * @Date 21:53 2022/3/18
 **/
func GetPostByID(pid int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id, title, content, author_id, author_name, community_id, status, create_time, like_num, 
	unlike_num, view_num, comment_num from post where post_id = ?`
	err = db.Get(post, sqlStr, pid)
	if err == sql.ErrNoRows {
		err = ErrorInvalidID
		return
	}
	if err != nil {
		zap.L().Error("query post failed", zap.String("sql", sqlStr), zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	return
}

/**
 * @Author wkwar
 * @Description //TODO 根据给定的id列表查询帖子数据
 * @Date 22:55 2022/3/18
 **/
func GetPostListByIDs(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, author_name, community_id, status, create_time, like_num, 
	unlike_num, view_num, comment_num from post 
	where post_id in (?)
	order by FIND_IN_SET(post_id, ?)`
	// 动态填充id
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return
	}
	// sqlx.In 返回带 `?` bindvar的查询语句, 我们使用Rebind()重新绑定它
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
	return
}

/**
 * @Author huchao
 * @Description //TODO 获取帖子列表 --- 根据创建时间获取
 * @Date 22:58 2022/2/12
 **/
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, author_name, community_id, status, create_time, like_num, 
	unlike_num, view_num, comment_num
	from post
	ORDER BY create_time
	DESC 
	limit ?,?
	`
	posts = make([]*models.Post, 0, 2)	// 0：长度  2：容量
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	return

}

func IncrCommentNum(id uint64) (err error) {
	sqlstr := "update post set comment_num = comment_num + 1 where post_id = ?"
	_, err = db.Exec(sqlstr, id)
	return
}


func GetPostsByRanking(num int) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, author_name, community_id, status, create_time, like_num, 
	unlike_num, view_num, comment_num from post ORDER BY view_num DESC limit ?`
	err = db.Select(&posts, sqlStr, num)
	return
}