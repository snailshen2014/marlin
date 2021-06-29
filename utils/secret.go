/*
 * @Author: David
 * @Date: 2020-03-09 21:58:30
 * @LastEditTime: 2020-03-09 21:58:38
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /sailfish/utils/secret.go
 */
package utils

import (
	"math/rand"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var n = 12

func GenerateSecret() string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(letterRunes))]
	}
	return string(b)
}
