package es

import (
	"fmt"
	"log"
	"os"
	"time"
	"context"
	"backbend/setting"
	"go.uber.org/zap"
	"github.com/olivere/elastic/v7"
)

//初始化 elastic

var (
	client *elastic.Client
	index  string
	ctx = context.Background()
)


const mappingTpl = `
{
	"mappings": {
		"properties": {
			"post_id": { 
				"type": "keyword" 
			},
			"author_id": { 
				"type": "keyword" 
			},
			"author_name": { 
				"type": "keyword" 
			},
			"community_id": {
				"type": "keyword" 
			},
			"status": {
				"type": "integer"
			},
			"title": { 
				"type": "text" 
			},
			"content":	{ 
				"type": "text" 
			},
			"create_time":	{ 
				"type": "date" 
			}
		}
	}
}`


type PostDocument struct {
	PostID      uint64    `json:"post_id,string"`
	AuthorId    uint64    `json:"author_id"`
	AuthorName  string 	  `json:"author_name" `
	CommunityID uint64     `json:"community_id"`
	Status      int32     `json:"status" `
	Title       string    `json:"title" `
	Content     string    `json:"content" `
	CreateTime  time.Time `json:"create_time"`
}

func Init(esConf *setting.ElasticConf) (err error) {
	index = esConf.Index
	// 创建ES client用于后续操作ES
	client, err = elastic.NewClient(
		// 设置ES服务地址，支持多个地址
		elastic.SetURL(esConf.Url),
		// 设置错误日志输出
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		// 设置info日志输出
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	if err != nil {
		// Handle error
		return
	}

	info, code, err := client.Ping(esConf.Url).Do(ctx)
    if err != nil {
        
        panic(err)
    }
	fmt.Println("info", info, "code", code)
	
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		return
	}

	if !exists {
		response, err := client.CreateIndex(index).BodyString(mappingTpl).Do(ctx)
		if err != nil {
			return err
		}
		zap.L().Info("创建elastic索引成功", zap.Any("index:", response.Index))
	}

	return
}
