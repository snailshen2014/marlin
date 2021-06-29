/*
 * @Author: David
 * @Date: 2020-03-08 12:19:03
 * @LastEditTime: 2020-11-14 19:38:20
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/datamodels/key.go
 */

package datamodels

//Key redis model
type Key struct {
	ID    string `json:"id" form:"id"`
	Value string `json:"value" form:"value"`
	TTL   int64  `json:"ttl" form:"ttl"`
}
