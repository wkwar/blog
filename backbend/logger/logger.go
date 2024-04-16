package logger

import (
	"os"
	"net"
	"time"
	"strings"
	"net/http"
	"runtime/debug"
	"net/http/httputil"
	"backbend/setting"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var lg *zap.Logger

//自定义日志输出格式
func Init(logConf *setting.LogConf, mode string) (err error) {
	//设置编码方式
	encoder := getEncoder()
	writeSyncer := getLogWriter(logConf.Filename, logConf.MaxSize, logConf.MaxBackups, logConf.MaxAge)
	//判断logConf里面的Level是否有问题 --- 是否是定义范围内
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(logConf.Level))
	if err != nil {
		return
	}

	var core zapcore.Core
	//判断logger模式，如果为dev模式就会输出到终端上
	if mode == "dev" {
		// 进入开发模式，日志输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		// 日志同时输出到控制台和日志文件中
		core = zapcore.NewTee(		// 多个输出
			zapcore.NewCore(encoder, writeSyncer, l),		// 往日志文件里面写
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),	// 终端输出
		)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}

	//创建日志记录器实例
	//zap.Addcaller() 输出日志打印文件和行数
	lg = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(lg)
	zap.L().Info("init logger success")
	return
}

//获取编码器，定义好写入日志的格式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder  // log 时间格式
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 输出level序列化为全大写
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)  // 以json格式写入
}

//定义日志保存位置 -- 包括写在哪个文件，其余参数
// 日志文件 与 日志切割 配置
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,	//日志文件格式
		MaxSize:    maxSize,	//单个日志文件最大 MB
		MaxBackups: maxBackup,	//日志备份数量
		MaxAge:     maxAge,		//日志最大保留时间
	}
	return zapcore.AddSync(lumberJackLogger)
}


// GinLogger 接收gin框架默认的日志
//相当于一个中间件（一个handler）
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		lg.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
//追踪 连接异常的原因，输出日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					lg.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}