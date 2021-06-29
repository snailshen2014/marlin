/*
 * @Author: David
 * @Date: 2020-03-08 20:59:18
 * @LastEditTime: 2020-11-14 19:37:41
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/conf/dbConfig.go
 */

package conf

import (
	"encoding/json"
	"io/ioutil"
)

type dbConfig struct {
	Port       string `json:"Port"`
	DBUserName string `json:"DBUserName"`
	DBPassword string `json:"DBPassword"`
	DBIp       string `json:"DBIp"`
	DBPort     string `json:"DBPort"`
	DBName     string `json:"DBName"`
}

var DbConfig = &dbConfig{}

func init() {
	//指定对应的json配置文件
	b, err := ioutil.ReadFile("dbConfig.json")
	if err != nil {
		panic("Sys config read err")
	}
	err = json.Unmarshal(b, DbConfig)
	if err != nil {
		panic(err)
	}

}
