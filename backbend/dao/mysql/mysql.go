package mysql
/**
mysql初始化 以及 关闭
**/

import (
	"fmt"
	"backbend/setting"
	"go.uber.org/zap"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sqlx.DB
)

//初始化数据库 db
func Init(cfg *setting.MySqlConf) (err error){
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true", cfg.User, cfg.Password, cfg.Host, cfg.Port,
		cfg.DB)
	
	db, err = sqlx.Open("mysql", dataSourceName)
	if err != nil {
		zap.S().Errorf("sql open database[%s] failed, err=%v", cfg.DB, err)
		return
	}
	//zap.S().Info("sql open database success")
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	return
}

func Close() {
	_ = db.Close()
}