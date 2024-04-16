package controller
//定义用户注册 以及 登录的路由方法
import (
	"fmt"
	"errors"
	"strings"
	"net/http"
	"backbend/models"
	"backbend/logic"
	"backbend/dao/mysql"

	"go.uber.org/zap"
	"github.com/gin-gonic/gin"
	"backbend/pkg/jwt"
	"github.com/go-playground/validator/v10"
)

// SignUpHandler 注册业务
// @Summary 注册业务
// @Description 注册业务
// @Tags 用户业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /signup [POST]
func SignUpHandler(c *gin.Context) {
	// 1.获取请求参数 2.校验数据有效性
	var fo = new(models.RegisterForm)
	if err := c.ShouldBindJSON(&fo); err != nil {
		//注册有问题，直接返回错误 
		zap.L().Error("signUp with invalid param", zap.Error(err))
		// 判断err是不是 validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseError(c, CodeInvalidParams)		// 请求参数错误
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		ResponseErrorWithMsg(c, CodeInvalidParams, removeTopStruct(errs.Translate(trans))) // 翻译错误
		return										
	}
	fmt.Println("fo = ", fo)
	// 3.业务处理——注册用户
	if err := logic.SignUp(fo); err != nil {
		zap.L().Error("logic.signup failed", zap.Error(err))
		//判断错误类型 
		if errors.Is(err, mysql.ErrorUserExit) {
			//用户名重复
			ResponseError(c, CodeUserExist)
			return
		}
		//其他错误，就返回服务器繁忙
		ResponseError(c, CodeServerBusy)
		return
	}
	zap.L().Info("注册用户成功")
	
	//4.返回响应,
	ResponseSuccess(c, nil)
}

// LoginHandler 登录业务
// @Summary 登录业务
// @Description 登录业务
// @Tags 用户业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.LoginForm false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /login [POST]
func LoginHandler(c *gin.Context) {
	// 1、获取请求参数及参数校验
	var u *models.LoginForm
	if err := c.ShouldBindJSON(&u); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param",zap.Error(err))
		// 判断err是不是 validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			ResponseError(c,CodeInvalidParams)		// 请求参数错误
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		ResponseErrorWithMsg(c, CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2、业务逻辑处理——登录
	user, err := logic.Login(u)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username",u.UserName), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExit) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidParams)
		return
	}
	// 3、返回响应 --- Token放在Header的Authorization中，并使用Bearer开头,这个不知道如何实现？
	fmt.Println("user", user)
	ResponseSuccess(c, gin.H{
		//存储用户信息，以及token信息
		"userInfo": gin.H{
			"user_id": fmt.Sprintf("%d", user.UserID), //js识别的最大值：id值大于1<<53-1  int64: i<<63-1
			"user_name": user.UserName,
		},
		"access_token": user.AccessToken,
		"refresh_token": user.RefreshToken,
	})
}

// RefreshTokenHandler 刷新accessToken
// @Summary 刷新accessToken
// @Description 刷新accessToken
// @Tags 用户业务接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /refresh_token [GET]
func RefreshTokenHandler(c *gin.Context) {
	//得到URL后面的参数
	rt := c.Query("refresh_token")
	// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
	// 这里假设Token放在Header的Authorization中，并使用Bearer开头
	// 这里的具体实现方式要依据你的实际业务情况决定
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		ResponseErrorWithMsg(c, CodeInvalidToken, "请求头缺少Auth Token")
		c.Abort()
		return
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		ResponseErrorWithMsg(c, CodeInvalidToken, "Token格式不对")
		c.Abort()
		return
	}
	//得到Token
	aToken, rToken, err := jwt.RefreshToken(parts[1], rt)
	fmt.Println(err)
	c.JSON(http.StatusOK, gin.H{
		"access_token":  aToken,
		"refresh_token": rToken,
	})
}