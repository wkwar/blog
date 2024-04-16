package models

import (
	"time"
)

/**
 * @Author wkwar
 * @Description //TODO 创建社区请求参数
 * @Date 14:00 2023/1/1
 **/
type Community struct {
	CommunityID   uint64 `json:"community_id" db:"community_id"`
	CommunityName string `json:"community_name" db:"community_name"`
}

/**
 * @Author wkwar
 * @Description //TODO 社区显示请求参数
 * @Date 14:00 2023/1/1
 **/
type CommunityDetail struct {
	CommunityID   uint64    `json:"community_id" db:"community_id"`
	CommunityName string    `json:"community_name" db:"community_name"`
	Introduction  string    `json:"introduction,omitempty" db:"introduction"`	// omitempty 当Introduction为空时不展示
	CreateTime    time.Time `json:"create_time" db:"create_time"`
}


