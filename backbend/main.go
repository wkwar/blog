package main

import (
	"fmt"
	"backbend/setting"
	"backbend/logger"
	"backbend/dao/es"
	"backbend/dao/mysql"
	"backbend/dao/redis"
	"backbend/pkg/check"
	"backbend/pkg/snowflake"
	"backbend/controller"
	"backbend/routers"
	
)

//@title bluebell_backend
//@version 1.0
//@description bluebell_backend测试
//@termsOfService http://swagger.io/terms/
//
//@contact.name author：@huchao
//@contact.url http://www.swagger.io/support
//@contact.email support@swagger.io
//
//@license.name Apache 2.0
//@license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
//@host 127.0.0.1:8081
//@BasePath /api/v1/
func main() {
	// 加载配置
	if err := setting.Init(); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}
	fmt.Println("config =", setting.Config)
	//日志输出初始化
	if err := logger.Init(setting.Config.LogConf, setting.Config.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	//mysql初始化
	if err := mysql.Init(setting.Config.MySqlConf); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close() // 程序退出关闭数据库连接
	//redis初始化
	if err := redis.Init(setting.Config.RedisConf); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()
	// 雪花算法生成分布式ID
	//起始时间格式为 xxxx-xx-xx
	if err := snowflake.Init("2023-01-01",1); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	//敏感词校验初始化
	if err := check.Init(setting.Config.Check.FilePath); err != nil {
		fmt.Printf("init check failed, err:%v\n", err)
		return
	}
	//es初始化
	if err := es.Init(setting.Config.ElasticConf); err != nil {
		fmt.Printf("init es failed, err:%v\n", err)
		return
	}


	//报错，可以使用中文输出
	if err := controller.InitValidator("zh");err!=nil{
		fmt.Printf("init validator Trans failed,err:%v\n",err)
		return
	}

	// 注册路由
	r := routers.SetupRouter(setting.Config.Mode)   
	err := r.Run(fmt.Sprintf(":%d", setting.Config.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}

