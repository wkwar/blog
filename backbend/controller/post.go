package controller

import (
	"backbend/dao/es"
	"backbend/logic"
	"backbend/models"
	"backbend/models/constants"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

//创建贴子处理
func CreatePostHandler(c *gin.Context) {
	// 1、获取参数及校验参数
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {   // validator --> binding tag
		zap.L().Debug("c.ShouldBindJSON(post) err",zap.Any("err",err))
		zap.L().Error("create post with invalid parm")
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}
	//根据得到的数据 post.ID 是否 == 0，区分为 创建 和 修改h'h'h'j'o'j'o
	fmt.Println("post", post)
	// 参数校验
	// 获取作者ID，当前请求的UserID(从c取到当前发请求的用户ID)
	user, err := getCurrentUser(c)
	fmt.Println("user", user)
	if err != nil {
		zap.L().Error("GetCurrentUser() failed", zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}
	post.AuthorId = user.UserID
	post.AuthorName = user.UserName
	post.Status = constants.StatusOk
	post.CreateTime = time.Now()
	//贴子校验
	if ok := CheckPost(&post); ok {
		post.Status = constants.StatusReview
	}
	
	//根据得到post
	if(post.PostID != 0) {
		//跳转到修改
		err = logic.UpdatePost(&post)
		if err != nil {
			zap.L().Error("logic.CreatePost failed", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}
		//更新es 的id数据
		err = es.UpdateIdPost(&post)
		if err != nil {
			zap.L().Error("es.UpdateIdPost failed", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}
	} else {
		// 2、创建帖子
		err = logic.CreatePost(&post)
		if err != nil {
			zap.L().Error("logic.CreatePost failed", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}

		//创建 es id数据
		err = es.CreatePostToEs(&post)
		if err != nil {
			zap.L().Error("es.CreatePostToEs failed", zap.Error(err))
			ResponseError(c, CodeServerBusy)
			return
		}
		
		// //生成帖子记录，包括 评论数，浏览量，点赞，收藏
		// err = logic.CreatePostRecord(int64(post.PostID))
		// if err != nil {
		// 	zap.L().Error("logic.CreatePostRecord failed", zap.Error(err))
		// 	ResponseError(c, CodeServerBusy)
		// 	return
		// }
		
	}

	// 3、返回响应  
	ResponseSuccess(c, nil)
}


//删除帖子内容
func DeletePostHandler(c *gin.Context) {
	//得到post_id
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {   // validator --> binding tag
		zap.L().Debug("c.ShouldBindJSON(post) err",zap.Any("err",err))
		zap.L().Error("create post with invalid parm")
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}

	err := logic.DeletePost(&post)
	if err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	err = es.DeleteIdPost(post.PostID)
	if err != nil {
		zap.L().Error("es.DeleteIdPost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3、返回响应  
	ResponseSuccess(c, nil)
}


//获得贴子列表 -- 分页展示
func PostListHandler(c *gin.Context) {
	// 获取分页参数 -- 页数以及当页的序号
	page, size := getPageInfo(c)
	// 获取数据
	data, err := logic.GetPostList(page, size)
	//fmt.Println("data", data)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	total, err := logic.GetTotalPost()
	if(err != nil) {
		ResponseError(c, CodeServerBusy)
		return
	}
	PostResponseSuccess(c, data, int64(total), page, size)
}


//升级版
func PostList2Handler(c *gin.Context)  {
	// GET请求参数(query string)： /api/v1/posts2?page=1&size=10&order=time
	// 获取分页参数
	p := &models.ParamPostList{
		Page: 1,
		Size: 10,
		Order: models.OrderTime,	// magic string
	}
	//c.ShouldBind() 根据请求的数据类型选择相应的方法去获取数据
	//c.ShouldBindJSON() 如果请求中携带的是json格式的数据，才能用这个方法获取到数据
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("PostList2Handler with invalid params",zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 获取数据
	data, err := logic.GetPostListNew(p)	// 更新：合二为一
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}


//根据id 获取帖子内容
func PostDetailHandler(c *gin.Context) {
	// 1、获取参数(从URL中获取帖子的id)
	postIdStr := c.Param("id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid param",zap.Error(err))
		ResponseError(c,CodeInvalidParams)
		return
	}

	fmt.Println("id", postIdStr)

	// 2、根据id取出id帖子数据(查缓存，数据库)
	post, err := logic.GetPostById(postId)
	if err != nil {
		zap.L().Error("logic.GetPost(postID) failed", zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}
	
	//浏览数 + 1
	//更新 帖子浏览数
	post.ViewNum += 1
	//数据库需要更新帖子
	fmt.Println("post", post.Post)
	err = logic.UpdatePost(post.Post)
	if err != nil {
		zap.L().Error("logic.UpdatePost(postID) failed", zap.Error(err))
		ResponseError(c,CodeServerBusy)
		return
	}

	// fmt.Println("记录XXXX")
	// record, err := logic.GetPostRecord(postId)
	// if err != nil {
	// 	zap.L().Error("logic.GetPostRecord failed", zap.Error(err))
	// 	ResponseError(c,CodeServerBusy)
	// }

	// err = record.ViewCount(postId)
	// if err != nil {
	// 	zap.L().Error("logic.GetPostRecord failed", zap.Error(err))
	// 	ResponseError(c,CodeServerBusy)
	// }

	
	// 3、返回响应
	ResponseSuccess(c, post)
}

/**
 * @Author huchao
 * @Description //TODO 根据社区去查询帖子列表
 * @Date 22:44 2022/2/16
 **/
func GetCommunityPostListHandler(c *gin.Context)  {
	// GET请求参数(query string)： /api/v1/posts2?page=1&size=10&order=time
	// 获取分页参数
	p := &models.ParamPostList{
		CommunityID: 0,
		Page:        1,
		Size:        10,
		Order:       models.OrderScore,
	}
	//c.ShouldBind() 根据请求的数据类型选择相应的方法去获取数据
	//c.ShouldBindJSON() 如果请求中携带的是json格式的数据，才能用这个方法获取到数据
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostListHandler with invalid params",zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 获取数据
	data, err := logic.GetCommunityPostList(p)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

func SearchPostHandler(c *gin.Context) {
	//根据 post请求，得到参数
	var search models.Search
	if err := c.ShouldBindJSON(&search); err != nil {   // validator --> binding tag
		zap.L().Error("c.ShouldBindJSON(search) err",zap.Any("err",err))
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}

	fmt.Println("search :", search)
	data, err := es.SearchPost(&search)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}


func GetSearchPostHandler(c *gin.Context) {
	//根据 post请求，得到参数
	var search models.Search
	// if err := c.ShouldBindJSON(&search); err != nil {   // validator --> binding tag
	// 	zap.L().Error("c.ShouldBindJSON(search) err",zap.Any("err",err))
	// 	ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
	// 	return
	// }
	search.KeyWord = c.Query("keyword")
	timeRange := c.Query("timeRange")
	search.TimeRange, _ = strconv.Atoi(timeRange)
	search.AuthorName = c.Query("author_name")

	fmt.Println("search :", search)
	data, err := es.SearchPost(&search)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

//根据数据 得到排行榜
func GetPostsByRanking(c *gin.Context) {
	//查询 top N 热度的博客
	data, err := logic.GetPostsByRanking(3)
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

