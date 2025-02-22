package controller

/**
 * @Author wkwar
 * @Description //TODO 自定义错误常量
 * @Date 14:00 2023/1/1
 **/
type MyCode int64

const (
	CodeSuccess         MyCode = 1000
	CodeInvalidParams   MyCode = 1001
	CodeUserExist       MyCode = 1002
	CodeUserNotExist    MyCode = 1003
	CodeInvalidPassword MyCode = 1004
	CodeServerBusy      MyCode = 1005

	CodeInvalidToken      MyCode = 1006
	CodeInvalidAuthFormat MyCode = 1007
	CodeNotLogin          MyCode = 1008
	CodeCreateFailed      MyCode = 1009
)

var msgFlags = map[MyCode]string{
	CodeSuccess:         "success",
	CodeInvalidParams:   "请求参数错误",
	CodeUserExist:       "用户名重复",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeInvalidToken:      "无效的Token",
	CodeInvalidAuthFormat: "认证格式有误",
	CodeNotLogin:          "未登录",
	CodeCreateFailed:      "创建数据失败",
}

//返回错误原因
func (c MyCode) Msg() string {
	msg, ok := msgFlags[c]
	if ok {
		return msg
	}
	return msgFlags[CodeServerBusy]
}
