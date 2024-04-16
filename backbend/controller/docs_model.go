package controller

import (
	"backbend/models"
)

// _ResponsePostList 帖子列表接口响应数据
type _ResponsePostList struct {
	Code    MyCode                  `json:"code"`    // 业务响应状态码
	Message string                  `json:"message"` // 提示信息
	Data    []*models.ApiPostDetail `json:"data"`    // 数据
}