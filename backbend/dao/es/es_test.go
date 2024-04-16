package es

import (
	"fmt"
	"time"
	"backbend/setting"
	"backbend/models"
	"testing"
)

func TestInsertData(t *testing.T) {
	var esConf = setting.ElasticConf{
		Url : "http://127.0.0.1:9200",
		Index : "test3",
	}
	err := Init(&esConf)
	if(err != nil) {
		t.Error("err:", err)
	}

	//判断插入贴子 数据有没有问题
	post := &models.Post{
		PostID: 5,
		AuthorId: 3,
		Title: "时间过期",
		Content: "测试的东西",
		CreateTime: time.Date(2023, 5, 6, 11, 45, 04, 0, time.UTC) ,
	}

	err = CreatePostToEs(post)
	if(err != nil) {
		t.Error("create error:", err)
	}

	// err = GetPostById(post.PostID)
	// if(err != nil) {
	// 	t.Error("get error:", err)
	// }

	//ESTermsQuery()
	serach := &models.Search{
		KeyWord: "无人机",
	}
	posts, err := SearchPost(serach)
	if(err != nil) {
		t.Error("get error:", err)
	}
	fmt.Println("posts", posts)
}