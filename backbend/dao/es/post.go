package es

import (
	"backbend/dao/mysql"
	"fmt"
	"time"
	"reflect"
	"encoding/json"
	"strconv"
	"backbend/models"
	"backbend/models/constants"
	"go.uber.org/zap"
	"github.com/olivere/elastic/v7"
	"github.com/mlogclub/simple/common/dates"

) 

func NewPostDoc(post *models.Post) (doc *PostDocument) {
	if post == nil {
		return nil
	}

	doc = &PostDocument{
		PostID      : post.PostID,
		AuthorId    : post.AuthorId,
		CommunityID : post.CommunityID,
		Status      : post.Status,
		Title       : post.Title,
		Content     : post.Content,
		CreateTime  : post.CreateTime,
	}

	// 处理内容
	// content := markdown.ToHTML(topic.Content)
	// content = html2.GetHtmlText(content)
	// content = html.EscapeString(content)

	// 处理用户
	user, _ := mysql.GetUserByID(post.AuthorId)
	
	if user != nil {
		doc.AuthorName = user.UserName
	}
	return doc

}
//插入 数据到 索引 中
func CreatePostToEs(post *models.Post) (err error) {
	doc := NewPostDoc(post)
	response, err := client.Index().
	Index(index). // 设置索引名称
	Id(strconv.FormatInt(int64(post.PostID), 10)). // 设置文档id
	BodyJson(doc). // 指定前面声明的微博内容
	Do(ctx)
	if err != nil {
		return
	}
	zap.L().Info("创建数据到es", zap.Any("result", response.Result))
	return
}

func GetPostById(id uint64) (err error) {
	// 根据id查询文档
	get1, err := client.Get().
	Index(index). // 指定索引名
	Id(strconv.FormatInt(int64(id), 10)). // 设置文档id
	Do(ctx) // 执行请求
	if err != nil {
		return
	}
	if get1.Found {
		zap.L().Info("根据索引，id查询es数据", zap.Any("result", get1.Source))
	}

	// 手动将文档内容转换成go struct对象
	msg2 := models.Post{}
	// 提取文档内容，原始类型是json数据
	data, _ := get1.Source.MarshalJSON()
	// 将json转成struct结果
	json.Unmarshal(data, &msg2)
	// 打印结果
	fmt.Println(msg2)
	return
}


func UpdateIdPost(post *models.Post) (err error) {
	doc := NewPostDoc(post)
	_, err = client.Update().
		Index(index). // 设置索引名
		Id(strconv.FormatInt(int64(post.PostID), 10)). // 文档id
		Doc(doc). 
		Do(ctx) // 执行ES查询
	if err != nil {
		return err
	}
	return
}

func DeleteIdPost(id uint64) (err error) {
	_, err = client.Delete().
			Index(index).
			Id(strconv.FormatInt(int64(id), 10)).
			Do(ctx)
	if err != nil {
		return
	}
	return
}


func ESTermsQuery() {
	matchQuery := elastic.NewMultiMatchQuery("测试", "content", "title")

	searchResult, err := client.Search().
		Index(index).   // 设置索引名
		Query(matchQuery).   // 设置查询条件
		Sort("create_time", true). // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
		From(0). // 设置分页参数 - 起始偏移量，从第0行记录开始
		Size(10).   // 设置分页参数 - 每页大小
		Do(ctx)             // 执行请求


	if err != nil {
		// Handle error
		panic(err)
	}
	
	fmt.Printf("查询消耗时间 %d ms, 结果总数: %d\n", searchResult.TookInMillis, searchResult.TotalHits())
	if searchResult.TotalHits() > 0 {
		// 查询结果不为空，则遍历结果
		var b1 models.Post
		// 通过Each方法，将es结果的json结构转换成struct对象
		for _, item := range searchResult.Each(reflect.TypeOf(b1)) {
			// 转换成Article对象
			if t, ok := item.(models.Post); ok {
				fmt.Println(t.Content)
			}
		}
	}
}



func SearchPost(s *models.Search) (posts []PostDocument, err error) {
	//贴子状态
	query := elastic.NewBoolQuery().
		Must(elastic.NewTermQuery("status", constants.StatusOk))

	if(s.AuthorName != "") {
		query = elastic.NewBoolQuery().
			Must(elastic.NewTermQuery("author_name", s.AuthorName))
	}

	//时间选择
	//范围查询
	if s.TimeRange == 1 { // 一天内
		beginTime := dates.Timestamp(time.Now().Add(-24 * time.Hour))
		query.Must(elastic.NewRangeQuery("create_time").Gte(beginTime))
	} else if s.TimeRange == 2 { // 一周内
		beginTime := dates.Timestamp(time.Now().Add(-7 * 24 * time.Hour))
		query.Must(elastic.NewRangeQuery("create_time").Gte(beginTime))
	} else if s.TimeRange == 3 { // 一月内
		beginTime := dates.Timestamp(time.Now().AddDate(0, -1, 0))
		query.Must(elastic.NewRangeQuery("create_time").Gte(beginTime))
	} else if s.TimeRange == 4 { // 一年内
		beginTime := dates.Timestamp(time.Now().AddDate(-1, 0, 0))
		query.Must(elastic.NewRangeQuery("create_time").Gte(beginTime))
	}

	query.Must(elastic.NewMultiMatchQuery(s.KeyWord, "title", "content"))

	highlight := elastic.NewHighlight().
		PreTags("<span class='search-highlight'>").PostTags("</span>").
		Fields(elastic.NewHighlighterField("title"), elastic.NewHighlighterField("content"))

	searchResult, err := client.Search().
		Index(index).
		Query(query).
		From(0).Size(10).
		Highlight(highlight).
		Do(ctx)
	if err != nil {
		return
	}
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	if totalHits := searchResult.TotalHits(); totalHits > 0 {
		var doc PostDocument
		// 通过Each方法，将es结果的json结构转换成struct对象
		for _, item := range searchResult.Each(reflect.TypeOf(doc)) {
			// 转换成Article对象
			posts = append(posts, item.(PostDocument))		
		}
		
	}
	return
}
