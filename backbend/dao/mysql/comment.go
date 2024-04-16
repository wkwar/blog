package mysql

import (
	"backbend/models"
	"fmt"

	"go.uber.org/zap"
)

func CreateComment(comment *models.Comment) (err error) {
	fmt.Println("comment", comment)
	sqlStr := `insert into comment(
	comment_id, content, post_id, author_id, author_name)
	values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, comment.CommentID, comment.Content, comment.PostID,
		comment.AuthorID, comment.AuthorName)
	if err != nil {
		zap.L().Error("insert comment failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

func GetCommentList(id int64) (commentList []*models.Comment, err error) {
	sqlStr := `select comment_id, content, post_id, author_id, author_name, create_time, reply_num, like_num 
	from comment where post_id = ?`
	
	err = db.Select(&commentList, sqlStr, id)
	if(err != nil) {
		return 
	}
	return
}


func IncrReplytNum(id uint64) (err error) {
	sqlstr := "update comment set reply_num = reply_num + 1 where comment_id = ?"
	_, err = db.Exec(sqlstr, id)
	return
}
