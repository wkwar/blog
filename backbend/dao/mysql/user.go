package mysql

import (
	"errors"
	"crypto/md5"
	"encoding/hex"
	"database/sql"

	"backbend/models"
)

const secret = "wkwar.vip"

/**
 * @Author wkwar
 * @Description //TODO 数据加密
 * @Date 14:00 2023/1/1
 **/
func encryptPassword(data []byte) (result string) {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum(data))
}

/**
 * @Author wkwar
 * @Description //TODO 核查用户是否存在数据库
 * @Date 14:00 2023/1/1
 **/
func CheckUserExist(username string) (err error) {
	sqlstr := `select count(user_id) from user where username = ?`
	var count int
	if err = db.Get(&count, sqlstr, username); err != nil {
		return 
	}
	if count > 0 {
		return errors.New("用户已存在")
	}
	return
}

/**
 * @Author wkwar
 * @Description //TODO 插入用户数据到数据库
 * @Date 14:00 2023/1/1
 **/
func InsertUser(user models.User) (err error) {
	// 对密码进行加密
	user.Password = encryptPassword([]byte(user.Password))
	// 执行SQL语句入库
	sqlstr := `insert into user(user_id,username,password) values(?,?,?)`
	_, err = db.Exec(sqlstr, user.UserID, user.UserName, user.Password)
	return 
}

/**
 * @Author wkwar
 * @Description //TODO 验证输入的登录信息
 * @Date 14:00 2023/1/1
 **/
func Login(user *models.User) (err error) {
	originPassword := user.Password // 记录一下原始密码(用户登录的密码)
	sqlStr := "select user_id, username, password from user where username = ?"
	err = db.Get(user, sqlStr, user.UserName)
	if err != nil && err != sql.ErrNoRows {
		// 查询数据库出错
		return
	}
	if err == sql.ErrNoRows {
		// 用户不存在
		return ErrorUserNotExit
	}
	// 生成加密密码与查询到的密码比较
	password := encryptPassword([]byte(originPassword))
	if user.Password != password {
		return ErrorPasswordWrong
	}
	return
}


/**
 * @Author wkwar
 * @Description //TODO 根据ID查询作者信息
 * @Date 22:05 2023/3/18
 **/
 func GetUserByID(id uint64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username from user where user_id = ?`
	err = db.Get(user, sqlStr, id)
	return
}
