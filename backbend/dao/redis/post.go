package redis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"backbend/models"

	"github.com/go-redis/redis/v8"
)

/**
 * @Author huchao
 * @Description //TODO 统计贴子数量
 * @Date 0:12 2022/2/17
 **/
 func UpdatePostsNums(num int) (int, error) {
	var key string = KeyPostsNums
	// n, err := client.Exists(ctx, key).Result()
	// if(n == 0) {
	// 	return 0, err
	// }

	value, _ := Get(key)
	v, _ := strconv.Atoi(value)
	if(num == 0) {
		return v, nil
	}
	
	v += num
	err := Set(key, v, 0)
	if(err != nil) {
		return -1, err
	}
	return v, nil
	
}

/**
 * @Author huchao
 * @Description //TODO 按照分数从大到小的顺序查询指定数量的元素
 * @Date 0:12 2022/2/17
 **/
 func getIDsFormKey(key string, page, size int64) ([]string, error) {
	start := (page-1) * size
	end := start + size - 1
	// 3.ZREVRANGE 按照分数从大到小的顺序查询指定数量的元素
	return client.ZRevRange(ctx, key, start, end).Result()
}

/**
 * @Author huchao
 * @Description //TODO 升级版投票列表接口：按创建时间排序 或者 按照 分数排序 (查询出的ids已经根据order从大到小排序)
 * @Date 22:19 2022/2/15
 **/
 func GetPostIDsInOrder(p *models.ParamPostList) ([] string, error)  {
	// 从redis获取id
	// 1.根据用户请求中携带的order参数确定要查询的redis key
	key := KeyPostTimeZSet		// 默认是时间
	if p.Order == models.OrderScore {	// 按照分数请求
		key = KeyPostScoreZSet
	}
	// 2.确定查询的索引起始点
	return getIDsFormKey(key, p.Page ,p.Size)
}

/**
 * @Author huchao
 * @Description //TODO 根据ids查询每篇帖子的投赞成票的数据
 * @Date 21:28 2022/2/16
 **/
func GetPostVoteData(ids []string) (data []int64, err error)  {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids{
	//	key := KeyPostVotedZSetPrefix + id
	//	// 查找key中分数是1的元素数量 -> 统计每篇帖子的赞成票的数量
	//	v := client.ZCount(key, "1", "1").Val()
	//	data = append(data, v)
	//}
	// 使用 pipeline一次发送多条命令减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids{
		key := KeyCommunityPostSetPrefix + id
		pipeline.ZCount(ctx, key, "1", "1")
	}
	cmders, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders{
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}

/**
 * @Author huchao
 * @Description //TODO 按社区查询ids(查询出的ids已经根据order从大到小排序)
 * @Date 23:06 2022/2/16
 * @Param orderKey:按照分数或时间排序
	将社区key与orderkey(社区或时间)做zinterstore
 **/
 
func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	// 1.根据用户请求中携带的order参数确定要查询的redis key
	orderkey := KeyPostTimeZSet		// 默认是时间
	if p.Order == models.OrderScore {	// 按照分数请求
		orderkey = KeyPostScoreZSet
	}

	// 使用zinterstore 把分区的帖子set与帖子分数的zset生成一个新的zset
	// 针对新的zset 按之前的逻辑取数据

	// 社区的key
	cKey := KeyCommunityPostSetPrefix + strconv.Itoa(int(p.CommunityID))

	// 利用缓存key减少zinterstore执行的次数 缓存key
	key := orderkey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(ctx, key).Val() < 1 {
		// 不存在，需要计算
		pipeline := client.Pipeline()
		//取交集，按最大的赋值
		pipeline.ZInterStore(ctx, key, &redis.ZStore{
			Keys: []string{cKey, orderkey},
			Aggregate: "MAX",	// 将两个zset函数聚合的时候 求最大值
		})		// zinterstore 计算
		pipeline.Expire(ctx, key, 60 * time.Second)	// 设置超时时间
		_, err := pipeline.Exec(ctx)
		if err != nil {
			return nil, err
		}
	}
	// 存在的就直接根据key查询ids
	return getIDsFormKey(key ,p.Page, p.Size)
}

// CreatePost 使用hash存储帖子信息
//zset 包含了 帖子下 每个人投票分数
func CreatePost(postID, userID uint64, title, summary string, CommunityID uint64) (err error) {
	now := float64(time.Now().Unix())
	votedKey := KeyPostVotedZSetPrefix + strconv.Itoa(int(postID))
	communityKey := KeyCommunityPostSetPrefix + strconv.Itoa(int(CommunityID))
	postInfo := map[string]interface{}{
		"title":    title,
		"summary":  summary,
		"post:id":  postID,
		"user:id":  userID,
		"time":     now,
		"like_num": 0,
		"unlike_num": 0,
		"view_num": 1,
		"collect_num": 0,
		"content_num": 0,
	}

	// 事务操作
	pipeline := client.TxPipeline()
	//有序集合 --- 都会关联一个double类型的分数score
	//帖子 按 投票分数 排序 
	pipeline.ZAdd(ctx, votedKey, &redis.Z{ // 作者默认投赞成票
		Score:  1,		//分数用来排序
		Member: userID,	//存储的value
	})
	//expire命令设置一个键的生存时间
	pipeline.Expire(ctx, votedKey, time.Second*OneWeekInSeconds) // 一周时间
	//HMSet 批量设置
	//使用 hash 存储 帖子的数据 
	pipeline.HMSet(ctx, KeyPostInfoHashPrefix + strconv.Itoa(int(postID)), postInfo)
	//zset 存储帖子分数排序
	pipeline.ZAdd(ctx, KeyPostScoreZSet, &redis.Z{ // 添加到分数的ZSet
		Score:  now + VoteScore,
		Member: postID,
	})
	pipeline.ZAdd(ctx, KeyPostTimeZSet, &redis.Z{ // 添加到时间的ZSet
		Score:  now,
		Member: postID,
	})
	//Set无序集合 添加
	pipeline.SAdd(ctx, communityKey, postID) // 添加到对应版块  把帖子添加到社区的set
	_, err = pipeline.Exec(ctx)
	return
}


//更新贴子信息  --- 这个有问题
// func UpdatePost(postID, userID uint64, title, summary string, CommunityID uint64) (err error) {
// 	now := float64(time.Now().Unix()) 
// 	postInfo := map[string]interface{}{
// 		"title":    title,
// 		"summary":  summary,
// 		"post:id":  postID,
// 		"user:id":  userID,
// 		"time":     now,
// 		"votes":    1,
// 		"comments": 0,
// 	}
// 	client.HMSet(ctx, KeyPostInfoHashPrefix + strconv.Itoa(int(postID)), postInfo)
// 	return
// }

//删除redis中的帖子信息
func DeletePost(postID uint64) (err error){
	//删除 postId 的贴子信息
	err = client.Del(ctx,  strconv.FormatInt(int64(postID), 10)).Err()
	if(err != nil) {
		return
	}
	// //删除 Zset里面的贴子
	// err = client.ZRem(ctx, KeyPostScoreZSet, postID).Err()
	// if(err != nil) {
	// 	return
	// }
	// err = client.ZRem(ctx, KeyPostTimeZSet, postID).Err()
	// if(err != nil) {
	// 	return
	// }
	return
} 


func LoadTopNPosts(data []*models.Post) (err error) {
	//将所有数据存储到 redis中
	//redis中使用hash表存储 post数据
	for _, post := range data {
		//如果存在 key
		key := strconv.FormatInt(int64(post.PostID), 10)
		if _, err := Get(key); err == nil {
			continue
		}
		//redis中创建hash
		data, _ := json.Marshal(post)
		err = Set(key, string(data), 0)
		if err != nil {
			return
		}
	}
	return 
}


func GetPostByID(postID int64) (post *models.Post, err error) {
	
	key := strconv.Itoa(int(postID))
	fmt.Println(key)
	data, err := Get(key)
	if err != nil {
		return 
	}
	if data == "0" {
		return
	}
	if err = json.Unmarshal([]byte(data), &post) ; err != nil {
		return
	}
	return
}
