package models

import ()

/**
 * @Author wkwar
 * @Description //TODO 用户参数
 * @Date 14:00 2023/1/1
 **/
type User struct {
	UserID   uint64 `json:"user_id,string" db:"user_id"` // 指定json序列化/反序列化时使用小写user_id
	UserName string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	AccessToken    string
	RefreshToken   string
}

/**
 * @Author wkwar
 * @Description //TODO 注册请求参数
 * @Date 14:00 2023/1/1
 **/
type RegisterForm struct {
	UserName        string `form:"RegisterForm.username" json:"username" binding:"required"`  //require代表必须要，
	Password        string `form:"RegisterForm.password" json:"password" binding:"required"`
	ConfirmPassword string `form:"RegisterForm.re_password" json:"re_password" binding:"required,eqfield=Password"`
}

/**
 * @Author wkwar
 * @Description //TODO 登录请求参数
 * @Date 14:00 2023/1/1
 **/
 type LoginForm struct {
	UserName        string `form:"username" json:"username" binding:"required"`
	Password        string `form:"password" json:"password" binding:"required"`
}

/**
 * @Author wkwar
 * @Description //TODO 投票请求参数
 * @Date 14:00 2023/1/1
 **/
type VoteDataForm struct {
	//UserID int 从请求中获取当前的用户
	PostID    string  `form:"post_id" json:"post_id" binding:"required"`	 // 帖子id
	Direction int8     `form:"direction" json:"direction,string" binding:"oneof=1 0 -1"`  // 赞成票(1)还是反对票(-1)取消投票(0)
}

