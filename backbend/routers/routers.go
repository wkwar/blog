package routers

import (
	"time"
	"net/http"
	"backbend/logger"
	"backbend/controller"
	"backbend/middlewares"

	_ "backbend/docs"  // 千万不要忘了导入把你上一步生成的docs

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	swaggoFiles "github.com/swaggo/files"
)

//路由注册
/**
 * @Author wkwar
 * @Description //TODO 设置路由
 * @Date 14:00 2022/3/17
 **/
func SetupRouter(mode string) *gin.Engine {
	//如果为发布模式
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	//添加中间件，request输出格式；panic recover
	r.Use(logger.GinLogger(), logger.GinRecovery(true),
		middlewares.RateLimitMiddleware(2 * time.Second, 1),
	)
	//跨域处理
	r.Use(middlewares.Cors())
	r.LoadHTMLFiles("backbend/templates/index.html")	// 加载html
	r.Static("/static", "./backbend/static")	// 加载静态文件
	//页面显示
	
	r.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})
	
	//注册swagger  --- 后台提供API 服务接口文档，是后端开发人员实现后台功能的接口以便提供给前端开发人员去实现界面功能
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggoFiles.Handler))

	//业务路由注册
	v1 := r.Group("/api/v1")
	v1.POST("/signup", controller.SignUpHandler)				//注册业务路由
	v1.POST("/login", controller.LoginHandler)					//登录业务
	v1.GET("/refresh_token", controller.RefreshTokenHandler)	//刷新Token
	//127.0.0.1:8081/api/v1/posts?page=1&size=1  一页多少个贴子
	v1.GET("/posts", controller.PostListHandler)		// 分页展示帖子列表
	v1.GET("/posts/rank", controller.GetPostsByRanking)
	v1.GET("/posts2", controller.PostList2Handler) 		// 根据时间或者分数排序分页展示帖子列表

	v1.GET("/community", controller.CommunityHandler)	// 获取分类社区列表
	v1.GET("/community/:id", controller.CommunityDetailHandler)	// 根据ID查找社区详情
	v1.GET("/post/:id", controller.PostDetailHandler) // 查询帖子详情

	//下面的路由 都是基于登陆成功后的操作，所以需要验证Token
	v1.Use(middlewares.JWTAuthMiddleware())	// 应用JWT认证中间件
	{
		v1.POST("/post", controller.CreatePostHandler)	 // 创建帖子
		v1.POST("/search", controller.SearchPostHandler)	 // 创建帖子
		v1.GET("/search", controller.GetSearchPostHandler)	 // 创建帖子
		v1.POST("/post/delete", controller.DeletePostHandler)	 // 删除帖子
		v1.POST("/vote", controller.VoteHandler)		   // 投票
		v1.POST("/comment", controller.CreateCommentHandler)
		v1.GET("/comment", controller.CommentListHandler)
		v1.POST("/reply", controller.CreateReplyHandler)
		v1.GET("/reply", controller.ReplyListHandler)
		v1.POST("/like/:id", controller.UserLikeHandler)
		//v1.GET("/comment/:postId", controller.CommentListByPostHandler)

		v1.GET("/ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})
	}

	return r
}