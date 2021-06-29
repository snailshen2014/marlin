/*
 * @Author: David
 * @Date: 2020-03-08 21:02:46
 * @LastEditTime: 2020-11-14 19:39:33
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/utils/deferutil.go
 */

package utils

import (
	"github.com/jinzhu/gorm"
)

func Defer(tx *gorm.DB, code *int) {
	if *code == 0 {
		//提交事务
		tx.Commit()
	} else {
		//回滚
		tx.Rollback()
	}
}
