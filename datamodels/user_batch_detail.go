/*
 * @Author: David
 * @Date: 2020-03-08 12:19:03
 * @LastEditTime: 2020-03-12 17:34:48
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/datamodels/key.go
 */

package datamodels

import "github.com/jinzhu/gorm"

//UserBatchDetail batch record
type UserBatchDetail struct {
	gorm.Model
	BatchID int64  `json:"batch_id"`
	KeyName string `json:"key_name"`
	Status  int64  `json:"status"`
}
