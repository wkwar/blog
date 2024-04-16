package logic

import (
	"backbend/dao/mysql"
	"backbend/models"
	"backbend/models/constants"
	"fmt"
)

func UserLike(u *models.UserLike) (err error) {
	//先判断 是否已经点赞
	ok := mysql.IsExists(int64(u.UserID), int64(u.PostID))
	if(ok) {
		//
		fmt.Println("已经点赞。。。")
		//点赞数 - 1，删除用户
		err = mysql.DecrPostLikeNum(int64(u.PostID))
		if(err != nil) {
			return
		}
		err = mysql.DeleteUserLike(u)
		if(err != nil) {
			return
		}
		//
		fmt.Println("删除成功")
		u.Status = constants.NotLikedYet
		return
	}

	//创建 到 数据库
	err = mysql.CreateUserLike(u)
	if(err != nil) {
		return
	}
	u.Status = constants.AlreadyLiked
	//点赞数 + 1
	err = mysql.IncrPostLikeNum(int64(u.PostID))
	if(err != nil) {
		return
	}
	return
}