package redis

import (
	"backbend/models"
	"fmt"
	"strconv"
)

func CreateReply(reply *models.Reply) (err error) {
	key := KeyReply + strconv.Itoa(int(reply.CommentID))
	fmt.Println("reply_key", key)
	//通过set集合存储 所有的评论数据
	err = client.LPush(ctx, key, reply).Err()
	if(err != nil) {
		return
	}
	
 	return
}

func GetReplyList(id uint64) (err error) {
	key := KeyReply + strconv.Itoa(int(id))
	data, err := client.Get(ctx, key).Result()
	if(err != nil) {
		return 
	}
	fmt.Println("data", data)
	return
}