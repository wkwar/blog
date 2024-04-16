package redis

//记录帖子 的 收藏量，点赞数，分享数，浏览量，

import (
	"strconv"
	// "encoding"
	// "encoding/json"
	
)

// var _ encoding.BinaryMarshaler = new(PostRecord)
// var _ encoding.BinaryUnmarshaler = new(PostRecord)

//redis 使用哈希表 存储这些数据
type PostRecord struct {
	LikeNum    int64 `json:"like_num" redis:"like_num"`
	UnlikeNum  int64 `json:"unlike_num" redis:"unlike_num"`
	ViewNum    int64 `json:"view_num" redis:"view_num"`
	CollectNum int64 `json:"collect_num" redis:"collect_num"`
	ContentNum int64 `json:"content_num" redis:"content_num"`
}


func CreatePostRecord(postId int64) (err error) {
	var mapHash = map[string]interface{} {
		"like_num": 0,
		"unlike_num": 0,
		"view_num": 1,
		"collect_num": 0,
		"content_num": 0,
	}
	err = HSetPostRecord(postId, mapHash)
	return
}
// func (p *PostRecord) MarshalBinary() (data []byte, err error) {
// 	return json.Marshal(p)
// }

// func (p *PostRecord) UnmarshalBinary(data []byte) (err error) {
// 	return json.Unmarshal(data, p)
// }

func (p *PostRecord) LikeCount(postId int64, v int64) (err error) {
	key := GetPostKey(postId)
	err = HIncrBy(key, "like_num", v)
	if err != nil {
		return
	}

	p.LikeNum += v
	return
}

func (p *PostRecord) UnLikeCount(postId int64, v int64) (err error) {
	key := GetPostKey(postId)
	err = HIncrBy(key, "unlike_num", v)
	if err != nil {
		return
	}

	p.UnlikeNum += v
	return
}

//浏览数只会增加
func (p *PostRecord) ViewCount(postId int64) (err error) {
	key := GetPostKey(postId)
	err = HIncrBy(key, "view_num", 1)
	if err != nil {
		return
	}
	p.ViewNum += 1
	return
}

func (p *PostRecord) CollectCount(postId int64, v int64) (err error) {
	key := GetPostKey(postId)
	err = HIncrBy(key, "collect_num", v)
	if err != nil {
		return
	}

	p.CollectNum += v
	return
}

func (p *PostRecord) ContentCount(postId int64, v int64) (err error) {
	key := GetPostKey(postId)
	err = HIncrBy(key, "content_num", v)
	if err != nil {
		return
	}

	p.ContentNum += v
	return
}

//简单的增加
func GetPostKey(postId int64) (key string) {
	return Prefix + "-" + strconv.Itoa(int(postId))
}


func Incr(key string) (err error){
	return client.Incr(ctx, key).Err()
}

func IncrBy(key string, value int64) (err error){
	return client.IncrBy(ctx, key,value).Err()
}

func HIncrBy(key string, filed string, value int64) (err error) {
	return client.HIncrBy(ctx, key, filed, value).Err()
}

func HSetPostRecord(postId int64, record interface{}) (err error) {
	key := GetPostKey(postId)
	err = client.HMSet(ctx, key, record).Err()
	return
}

func HGetPostRecord(postId int64) (record PostRecord,err error) {
	key := GetPostKey(postId)
	err = client.HGetAll(ctx, key).Scan(&record)
	return
}