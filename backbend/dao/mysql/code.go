package mysql

import (
	"errors"
)
/**
 * @Author wkwar
 * @Description //TODO 自定义数据库类型的错误常量
 * @Date 14:00 2023/1/1
 **/

var (
	ErrorUserExit      = errors.New("用户已存在")
	ErrorUserNotExit   = errors.New("用户不已存在")
	ErrorPasswordWrong = errors.New("密码错误")
	ErrorGenIDFailed   = errors.New("创建用户ID失败")
	ErrorInvalidID     = errors.New("无效的ID")
	ErrorQueryFailed   = errors.New("查询数据失败")
	ErrorInsertFailed  = errors.New("插入数据失败")
)
