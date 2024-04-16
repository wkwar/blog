package logic

import "backbend/dao/redis"



func CreatePostRecord(postId int64) (err error) {
	err = redis.CreatePostRecord(postId)
	return
}

func GetPostRecord(postId int64) (redis.PostRecord, error) {
	return redis.HGetPostRecord(postId)
}