package controller

import (
	"backbend/logic"
	"backbend/models"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserIDKey = "userID"
)

var (
	ErrorUserNotLogin = errors.New("当前用户未登录")
)


func getPageInfo(c *gin.Context) (int64, int64)  {
	pageStr := c.Query("page")
	SizeStr := c.Query("size")

	var (
		page int64		// 第几页 页数
		size int64	    // 每页几条数据
		err error
	)
	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(SizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}

//从request请求中得到用户ID
func getCurrentUser(c *gin.Context) (user *models.User, err error) {
	_userID, ok := c.Get(ContextUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok := _userID.(uint64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}

	//根据用户id 得到 用户信息
	user, err = logic.GetUserByID(userID)
	return
}