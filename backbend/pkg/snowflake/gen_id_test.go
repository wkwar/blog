package snowflake

import (
	"testing"
)

func TestGenID(t *testing.T) {
	//起始时间 必须为 xxxx-xx-xx 格式
	err := Init("2022-03-15", 1)
	if err != nil {
		t.Fatal("sonyflake init failed", err)
		return
	}
	id, err := sonyFlake.NextID()
	if err != nil {
		t.Fatal("get id failed", err)
		return
	}
	t.Logf("get id = %d", id)
}