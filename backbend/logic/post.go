package logic

import (
	"fmt"
	"backbend/models"
	"backbend/dao/mysql"
	"backbend/pkg/snowflake"
	"backbend/dao/redis"
	"go.uber.org/zap"
)

/**
 * @Author huchao
 * @Description //TODO 创建贴子
 * @Date 21:39 2022/2/12
 **/
func CreatePost(post *models.Post) (err error) {
	// 1、 生成post_id(生成帖子ID)
	postID, err := snowflake.GetID()
	if err != nil {
		zap.L().Error("snowflake.GetID() failed", zap.Error(err))
		return
	}
	post.PostID = postID
	// 2、创建帖子 保存到数据库
	if err := mysql.CreatePost(post); err != nil {
		zap.L().Error("mysql.CreatePost(&post) failed", zap.Error(err))
		return err
	}
	// community, err := mysql.GetCommunityNameByID(fmt.Sprint(post.CommunityID))
	// if err != nil {
	// 	zap.L().Error("mysql.GetCommunityNameByID failed", zap.Error(err))
	// 	return err
	// }
	// 3、redis存储帖子信息
	// if err := redis.CreatePost(
	// 	post.PostID,
	// 	post.AuthorId,
	// 	post.Title,
	// 	TruncateByWords(post.Content, 120),	//信息梗概
	// 	community.CommunityID,
	// 	); err != nil {
	// 	zap.L().Error("redis.CreatePost failed", zap.Error(err))
	// 	return err
	// }
	//更新贴子数量到redis中
	redis.UpdatePostsNums(1)
	return

}

/**
 * @Author huchao
 * @Description //TODO 更新贴子
 * @Date 21:39 2022/2/12
 **/

//先写数据库，后删除缓存
func UpdatePost(post *models.Post) (err error) {
	// 1、数据库更新帖子
	if err = mysql.UpdatePost(post); err != nil {
		zap.L().Error("mysql.CreatePost(&post) failed", zap.Error(err))
		return err
	}
	
	//2、删除redis中帖子信息
	if err = redis.DeletePost(post.PostID); err != nil {
		zap.L().Error("redis.DeletePost(&post) failed", zap.Error(err))
		return err
	}
	
	return
}


func DeletePost(post *models.Post) (err error) {
	// 2、更新帖子 保存到数据库
	postId := post.PostID
	if err := mysql.DeletePost(postId); err != nil {
		zap.L().Error("mysql.DeletePost(&post) failed", zap.Error(err))
		return err
	}
	
	if err := redis.DeletePost(postId); err != nil {
		zap.L().Error("redis.DeletePost(&post) failed", zap.Error(err))
		return err
	}
	//更新贴子数量到redis中
	redis.UpdatePostsNums(-1)
	return
}
/**
 * @Author huchao
 * @Description //TODO 得到贴子总数
 * @Date 21:39 2022/2/12
 **/
func GetTotalPost() (int, error) {
	return redis.UpdatePostsNums(0)
}


/**
 * @Author huchao
 * @Description //TODO 根据Id查询帖子详情
 * @Date 21:39 2022/2/12
 **/
 func GetPostById(postID int64) (data *models.ApiPostDetail, err error) {
	// 查询并组合我们接口想用的数据
	//先从redis缓存中查找数据
	fmt.Println("xxxxx")
	post, err := redis.GetPostByID(postID) 
	if err != nil {
		// zap.L().Debug("redis.GetPostByID(postID) failed",
		// 	zap.Int64("postID",postID),
		// 	zap.Error(err))
		
		fmt.Println("从数据库中查询数据")
		//从数据库中查询帖子信息
		post, err = mysql.GetPostByID(postID)
		if err != nil {
			zap.L().Error("mysql.GetPostByID(postID) failed",
				zap.Int64("postID",postID),
				zap.Error(err))
			return nil, err
		}
	} else {
		fmt.Println("从缓存Redis中查询数据")
	}
	
	// 根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorId)
	if err != nil {
		zap.L().Error("mysql.GetUserByID() failed",
			zap.Uint64("postID",post.AuthorId),
			zap.Error(err))
		return
	}

	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID() failed",
			zap.Uint64("community_id",post.CommunityID),
			zap.Error(err))
		return
	}
	// 接口数据拼接
	data = &models.ApiPostDetail{
		Post:            post,
		CommunityDetail: community,
		AuthorName:      user.UserName,
	}

	return
}

//获取贴子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	postList, err := mysql.GetPostList(page, size)
	if err != nil {
		fmt.Println(err)
		return
	}
	data = make([]*models.ApiPostDetail, 0, len(postList))	// data 初始化
	for _, post := range postList {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Uint64("postID",post.AuthorId),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed",
				zap.Uint64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		// 接口数据拼接
		postdetail := &models.ApiPostDetail{
			Post:            post,
			CommunityDetail: community,
			AuthorName:      user.UserName,
			Page: page,
			Size: size,
		}
		data = append(data, postdetail)
	}
	return
}

/**
 * @Author huchao
 * @Description //TODO 将两个查询帖子列表逻辑合二为一的函数
 * @Date 12:08 2022/2/17
 **/
 func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 根据请求参数的不同,执行不同的业务逻辑
	if p.CommunityID == 0 {
		// 查所有
		data, err = GetPostList2(p)
	} else {
		// 根据社区id查询
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed",zap.Error(err))
		return nil, err
	}
	return
}

/**
 * @Author huchao
 * @Description //TODO 升级版帖子列表接口：按创建时间排序 或者 按照 分数排序
 * @Date 22:03 2022/2/15
 **/
 func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2、去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 3、根据id去数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id的顺序返回  order by FIND_IN_SET(post_id, ?)
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Uint64("postID",post.AuthorId),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed",
				zap.Uint64("community_id",post.CommunityID),
				zap.Error(err))
			continue
		}
		// 接口数据拼接
		postdetail := &models.ApiPostDetail{
			VoteNum: voteData[idx],
			Post:            post,
			CommunityDetail: community,
			AuthorName:      user.UserName,
		}
		data = append(data,postdetail)
	}
	return
}

/**
 * @Author wkwar
 * @Description //TODO  根据社区去查询帖子列表
 * @Date 22:53 2022/2/16
 **/
func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 2、去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostList(p) return 0 data")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids", ids))
	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 3、根据id去数据库查询帖子详细信息
	// 返回的数据还要按照我给定的id的顺序返回  order by FIND_IN_SET(post_id, ?)
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserByID() failed",
				zap.Uint64("postID",post.AuthorId),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID() failed",
				zap.Uint64("community_id",post.CommunityID),
				zap.Error(err))
			continue
		}
		// 接口数据拼接
		postdetail := &models.ApiPostDetail{
			VoteNum: voteData[idx],
			Post:            post,
			CommunityDetail: community,
			AuthorName:      user.UserName,
		}
		data = append(data,postdetail)
	}
	return
}


func GetPostsByRanking(num int) (data []*models.Post, err error) {
	//从数据库中获取排名
	data, err = mysql.GetPostsByRanking(num)
	//需要将这些数据加载到缓存中
	err = redis.LoadTopNPosts(data)
	return
}