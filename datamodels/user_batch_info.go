/*
 * @Author: David
 * @Date: 2020-03-08 12:19:03
 * @LastEditTime: 2020-11-14 19:38:28
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/datamodels/key.go
 */

package datamodels

import "github.com/jinzhu/gorm"

//UserBatchInfo batch delete
type UserBatchInfo struct {
	gorm.Model
	UserName string `json:"user_name"`
	FileName string `json:"file_name"`
	Status   int64  `json:"status"`
}
