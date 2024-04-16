package jwt

import (
	//"fmt"
	"testing"
)

/**
 * @Author wkwar
 * @Description //TODO 自定义错误常量
 * @Date 14:00 2023/1/1
 **/
func TestJwtToken(t *testing.T) {
	aToken, rToken, err := GenToken(1, "123")
	if err != nil {
		t.Fatal("genToken failed,", err)
		return
	}
	t.Logf("aToken=%s, rToken=%s", aToken, rToken)
	//验证解析Token
	acliam, err := ParseToken(aToken)
	if err != nil {
		t.Fatal("genToken failed,", err)
		return
	}
	t.Log("parse Token cliam=", acliam)

	rcliam, err := ParseToken(rToken)
	if err != nil {
		t.Fatal("genToken failed,", err)
		return
	}
	t.Log("parse Token cliam=", rcliam)

	//验证刷新Token
	newaToken, newrToken, err :=  RefreshToken(aToken, rToken)
	if err != nil {
		t.Fatal("genToken failed,", err)
		return
	}
	t.Logf("newaToken=%s, newrToken=%s", newaToken, newrToken)
}