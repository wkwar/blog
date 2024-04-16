package redis

import (
	"backbend/models"
	"fmt"
	"strconv"
)

func CreateComment(comment *models.Comment) (err error) {
	key := KeyCommentList + strconv.Itoa(int(comment.PostID))

	//通过set集合存储 所有的评论数据
	err = client.SAdd(ctx, key, comment).Err()
	if(err != nil) {
		return
	}
	
 	return
}

func GetCommentList(id uint64) (err error) {
	key := KeyCommentList + strconv.Itoa(int(id))
	fmt.Println("comment_key", key)
	data, err := client.Get(ctx, key).Result()
	if(err != nil) {
		return 
	}
	fmt.Println("data", data)
	return
}