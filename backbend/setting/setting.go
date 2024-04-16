package setting

import (
	"github.com/spf13/viper"
)

/*
将配置文件加载初始化
*/

var Config = new(AppConfig)

//总配置文件，包含了 mysql,redis,logger配置，还有连接
type AppConfig struct {
	Mode         string `mapstructure:"mode"`
	Port         int    `mapstructure:"port"`
	Name         string `mapstructure:"name"`
	Version      string `mapstructure:"version"`
	StartTime    string `mapstructure:"start_time"`
	MachineID    int    `mapstructure:"machine_id"`
	*LogConf     `mapstructure:"log"`
	*MySqlConf   `mapstructure:"mysql"`
	*RedisConf   `mapstructure:"redis"`
	*Check       `mapstructure:"check"`
	*ElasticConf `mapstructure:"elastic"`
}

//mysql 配置
type MySqlConf struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

//redis配置
type RedisConf struct {
	Host         string `mapstructure:"host"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type LogConf struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type Check struct {
	FilePath string `mapstructure:"file_path"`
}

type ElasticConf struct {
	Url   string `mapstructure:"url"`
	Index string `mapstructure:"index"`
}

//读取mapstructure文件 有两种方式，一种是mapstructure.v2 读取
//一种使用 viper读取
func Init() (err error) {
	//加载配置文件
	//viper.SetConfigType("yaml")
	viper.SetConfigFile("./conf/config.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	//映射到结构体中

	err = viper.Unmarshal(&Config)
	if err != nil {
		return
	}
	return
}
