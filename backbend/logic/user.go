package logic

import (
	"backbend/models"
	"backbend/dao/mysql"
	"backbend/pkg/snowflake"
	"backbend/pkg/jwt"
)


func SignUp(p *models.RegisterForm) (err error) {
	// 1、判断用户存不存在
	err = mysql.CheckUserExist(p.UserName)
	if err != nil {
		// 数据库查询出错 或者 已经存在该用户
		return 
	}

	// 2、生成UID
	userId, err := snowflake.GetID()
	if err != nil {
		return mysql.ErrorGenIDFailed
	}
	// 构造一个User实例
	u := models.User{
		UserID:   userId,
		UserName: p.UserName,
		Password: p.Password,
	}
	// 3、保存进数据库
	return mysql.InsertUser(u)
}


func Login(p *models.LoginForm) (user *models.User, err error) {
	user = &models.User{
		UserName: p.UserName,
		Password: p.Password,
	}
	if err = mysql.Login(user); err != nil {
		return 
	}
	// 登录成功后，使用JWT生成Token -- Token将登录的信息进行加密，
	atoken, rtoken, err := jwt.GenToken(user.UserID, user.UserName)
	if err != nil {
		return
	}
	user.AccessToken = atoken
	user.RefreshToken = rtoken
	return
}


func GetUserByID(id uint64) (*models.User, error) {
	return mysql.GetUserByID(id)
}