package controller

import (
	"testing"
	"net/http"
	"github.com/gin-gonic/gin"
 	"github.com/go-playground/validator/v10"
	
)

/**
 * @Author wkwar
 * @Description //TODO 测试 翻译方法
 * @Date 14:00 2022/3/17
 **/

//翻译使用方法
// 1.errs, ok := err.(validator.ValidationErrors)  转成ValidationErrors形式
// 2.errs.Translate(Trans) 	翻译

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,max=16,min=6"`
}

func TestValidate(t *testing.T) {
	var v *validator.Validate
	if err := InitTrans(v, "zh"); err != nil {
		t.Fatalf("init trans failed, err:%v\n", err)
		return
	}
	router := gin.Default()
	router.POST("/user/login", login)
	router.Run(":8888")	
}
   
func login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 获取validator.ValidationErrors类型的errors
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非validator.ValidationErrors类型错误直接返回
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		// validator.ValidationErrors类型错误则进行翻译
		c.JSON(http.StatusOK, gin.H{
			"msg": errs.Translate(trans),
		})
		return
	}
	//login 操作省略
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}