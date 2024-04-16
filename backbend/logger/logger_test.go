package logger

import (
	"testing"
	"go.uber.org/zap"
	"backbend/setting"
)

/*
测试自定义logger
*/


var logConf = &setting.LogConf {
	Level: "debug",
	Filename:       "./test.log", // 输出日志文件名称
	MaxSize:    1,          // 输出单个日志文件大小，单位MB
	MaxBackups: 10,         // 输出最大日志备份个数
	MaxAge:         1000,       // 日志保留时间，单位: 天 (day)
}

func TestLogger(t *testing.T) {
	//输入dev模式，终端也会显示结果
	//否则只显示在文件中 --- 该文件是 该代码下对应的目录开始
	if err := Init(logConf, "dev"); err != nil {
		t.Fatal(err) 
	}
	//logger输出格式
	zap.S().Infof("测试，Infof模式:%s", "111")
	zap.S().Debug("测试Debud模式", "sucess")

}