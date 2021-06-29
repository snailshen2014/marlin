/*
 * @Author: David
 * @Date: 2020-03-08 21:00:59
 * @LastEditTime: 2020-11-14 19:38:40
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/datasource/datasource.go
 */

package datasource

import (
	"fmt"
	"marlin/conf"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func GetDB() *gorm.DB {
	return db
}

func init() {
	fmt.Println("datasource init running...")
	path := strings.Join([]string{conf.DbConfig.DBUserName, ":", conf.DbConfig.DBPassword, "@(", conf.DbConfig.DBIp, ":", conf.DbConfig.DBPort, ")/", conf.DbConfig.DBName, "?charset=utf8&parseTime=true"}, "")
	var err error
	db, err = gorm.Open("mysql", path)
	if err != nil {
		panic(err)
	}
	db.SingularTable(true)
	db.DB().SetConnMaxLifetime(1 * time.Second)
	db.DB().SetMaxIdleConns(20) //最大打开的连接数
	db.DB().SetMaxOpenConns(30) //设置最大闲置个数
	db.SingularTable(true)      //表生成结尾不带s
	// 启用Logger，显示详细日志
	db.LogMode(true)
	Createtable()
}
