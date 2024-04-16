package redis

import (
	//"chitchat/backbend/models"
	//"encoding/json"
	"fmt"
	"testing"
	//"time"
)

type post struct {
	id int
	name string
}

func TestRedis(t *testing.T) {
	// post := models.Post {
	// 	PostID: 1,
	// 	AuthorName: "小明",
	// }
	p := post{
		id : 1,
		name: "tom",
	}
	key := "test"
	
	//vals, err := json.Marshal(post)
	// fmt.Println(key, string(vals))
	// if err != nil {
	// 	t.Error(err)
	// }
	fmt.Println("xxx", p)
	err := client.Set(ctx, key, "aaa", 0).Err()
	fmt.Println("xxx")
	//err = Set(key, string(vals), time.Minute)
	if err != nil {
		
		t.Error("err=", err)
	}
	fmt.Println("xxx")
	// val, err := Get(key)
	// fmt.Println(val)
	// if err != nil {
	// 	t.Error(err)
	// }
}