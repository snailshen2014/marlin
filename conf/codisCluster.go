/*
 * @Author: David
 * @Date: 2020-03-08 20:59:18
 * @LastEditTime: 2020-03-08 20:59:19
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/conf/dbConfig.go
 */

package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type clusterConfig struct {
	ClusterName string `json:"cluster_name"`
	ZkAddr      string `json:"zk_addr"`
}

var clusters []clusterConfig

func init() {
	//指定对应的json配置文件
	b, err := ioutil.ReadFile("cluster.json")
	if err != nil {
		panic("Cluster config read err")
	}
	err = json.Unmarshal(b, &clusters)
	if err != nil {
		panic(err)
	}

}

//GetClusterAddr from config
func GetClusterAddr(clusterName string) string {
	for _, value := range clusters {
		fmt.Printf("value:%s\n", value)
		if value.ClusterName == clusterName {
			return value.ZkAddr
		}
	}
	return ""
}
