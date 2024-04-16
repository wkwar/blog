package mysql

import (
	"backbend/models"
	"fmt"

	"go.uber.org/zap"
)

func IsExists(userId, postId int64) bool {
	sqlstr := "select create_time from likes where user_id = ? and post_id = ?"
	var data = new(models.UserLike)
	err := db.Get(data, sqlstr, userId, postId)
	

	fmt.Println("是否点赞", data, err)
	if err != nil {
		return false
	}
	
	return true
}

func CreateUserLike(u *models.UserLike) (err error) {
	sqlstr := "insert into likes(user_id, post_id, create_time) values(?, ?, ?)"
	_, err = db.Exec(sqlstr, u.UserID, u.PostID, u.CreateTime)
	if err != nil {
		zap.L().Error("insert likes failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}


func DeleteUserLike(u *models.UserLike) (err error) {
	sqlstr := "delete from likes where user_id = ? and post_id = ?"
	_, err = db.Exec(sqlstr, u.UserID, u.PostID)
	if err != nil {
		zap.L().Error("delete likes failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

func IncrPostLikeNum(postId int64) (err error) {
	sqlstr := "update post set like_num = like_num + 1 where post_id = ?"
	_, err = db.Exec(sqlstr, postId)
	return
}

func DecrPostLikeNum(postId int64) (err error) {
	sqlstr := "update post set like_num = like_num - 1 where post_id = ?"
	_, err = db.Exec(sqlstr, postId)
	return
}