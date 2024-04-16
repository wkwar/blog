package logic

import (
	"backbend/models"
	"backbend/dao/mysql"
)


/**
 * @Author huchao
 * @Description //TODO 查询分类社区列表
 * @Date 16:42 2022/2/12
 **/
 func GetCommunityList() ([] *models.Community,error) {
	// 查数据库 查找到所有的community 并返回
	return mysql.GetCommunityList()
}

/**
 * @Author huchao
 * @Description //TODO 根据ID查询分类社区详情 
 * @Date 17:08 2022/2/12
 **/
func GetCommunityDetailByID(id uint64) (*models.CommunityDetail,error) {
	return mysql.GetCommunityByID(id)
}
