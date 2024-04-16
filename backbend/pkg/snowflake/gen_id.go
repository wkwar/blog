package snowflake
import (
	"fmt"
	"time"
	"github.com/sony/sonyflake"
)
/**
使用雪花算法来生成唯一，递增ID(64bit)  == 没有使用(1bit) + 时间戳(41bit) + 机器ID(10bit) + 序列号(12bit)

获取当前的毫秒时间戳；
用当前的毫秒时间戳和上次保存的时间戳进行比较；
如果和上次保存的时间戳相等，那么对序列号 sequence 加一；
如果不相等，那么直接设置 sequence 为 0 即可；

然后通过或运算拼接雪花算法需要返回的 int64 返回值。
**/



var (
	sonyFlake *sonyflake.Sonyflake 	//实例
	machineID uint16				//机器ID
)

//生成雪花ID，需要找到起始时间--- 时间戳， 机器ID -- 自己定义的，序列号 --- 自己增加（根据时间戳）
func Init(startTime string, machinID int64) (err error) {
	// 格式化 1月2号下午3时4分5秒  2006年
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	//
	settings := sonyflake.Settings{
		StartTime:      st,
		MachineID:     getMachineID,
		//CheckMachineID func(uint16) bool
	}
	//按照设置生成实例
	sonyFlake = sonyflake.NewSonyflake(settings)
	return
	
}

//满足setting结构体，所创建的一个函数
func getMachineID() (uint16, error) {
	return machineID, nil
}

// 生成 64 位的 雪花 ID
func GetID() (id uint64, err error) {
	if sonyFlake == nil {
		err = fmt.Errorf("sonyflake has not init")
		return 
	}
	return sonyFlake.NextID()
}
