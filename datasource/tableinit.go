/*
 * @Author: David
 * @Date: 2020-03-08 21:01:22
 * @LastEditTime: 2020-03-08 21:35:06
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/datasource/tableinit.go
 */
package datasource

import (
	"marlin/datamodels"
)

// Createtable 初始化表 如果不存在该表 则自动创建
func Createtable() {
	GetDB().AutoMigrate(
		&datamodels.User{},
		&datamodels.UserBatchInfo{},
		&datamodels.UserBatchDetail{},
		&datamodels.AppAccessAuthorication{},
	)
}
