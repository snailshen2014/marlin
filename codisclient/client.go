/*
 * @Author: David
 * @Date: 2020-03-09 23:33:24
 * @LastEditTime: 2020-11-14 19:37:49
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /sailfish/codis/client.go
 */

package codisclient

import (
	"fmt"
	"marlin/conf"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

var pool *Pool

var zkAddr string
var productName string

//SetZkInfo ,set productName and addr
func SetZkInfo(pName string) {
	productName = pName
	zkAddr = conf.GetClusterAddr(pName)
	initPool()
}

//GetPool ,for get codis pool
func GetPool() *Pool {
	return pool
}

//Close connection
func Close() {
	GetPool().Close()
}

func initPool() {
	fmt.Printf("Zk info ,addr:%s,zkdir:%s\n", zkAddr, "/codis3/"+productName+"/proxy")
	zkdir := "/codis3/" + productName + "/proxy"
	pool = &Pool{
		ZkServers: strings.Split(zkAddr, ","),
		ZkTimeout: time.Second * 60,
		ZkDir:     zkdir,
		Dial: func(network, address string) (redis.Conn, error) {
			conn, err := redis.Dial(network, address)
			if err != nil {
				conn.Send("AUTH", "PASSWORD")
			}
			return conn, err
		},
	}
}

//SscanSlots for slotsscan all slots for codis
func SscanSlots(from string, keys chan string) {
	client, err := pool.Get()
	if err != nil {
		return
	}
	defer client.Close()

	allSlotSum := 0
	for slot := 0; slot < 1024; slot++ {
		start := time.Now()
		cursor := 0
		slotNum := 0
		for {
			resultValue, err := redis.Values(client.Do("SLOTSSCAN", slot, cursor))
			if err != nil {
				fmt.Println(err)
				continue
			}
			slotKeys, ok := redis.Strings(resultValue[1], nil)
			if ok != nil {
				fmt.Printf("SLOTSSCAN slot:%d,cursor:%d,error:%s\n", slot, cursor, ok)
				continue
			}
			slotNum += len(slotKeys)
			for _, key := range slotKeys {
				if strings.HasPrefix(key, from) {
					keys <- key
				}

			}

			cursor, _ = redis.Int(resultValue[0], nil)
			if cursor == 0 {
				cost := time.Since(start)
				fmt.Printf("######slot:%d sscan finished,keys size:%d,time cost:%v\n",
					slot, slotNum, cost)

				allSlotSum += slotNum
				break
			}

		}

	}
	fmt.Printf("######all slots sscan finished,keys size:[%d]\n", allSlotSum)
	close(keys) //close channel
}

//KeyType , return key type
func KeyType(key string) (string, error) {
	client, err := pool.Get()
	if err != nil {
		return "", err
	}
	defer client.Close()
	result, err := redis.String(client.Do("type", key))
	if err == nil {
		return result, nil
	}
	return "", err
}

//Get get value
func Get(productName, key string) (string, error) {
	SetZkInfo(productName)
	client, err := pool.Get()
	if err != nil {
		return "zk", err
	}
	defer client.Close()
	keyType, err := KeyType(key)
	if err != nil {
		return "", err
	}
	if keyType == "none" {
		return keyType, nil
	}
	switch keyType {
	case "string":
		return redis.String(client.Do("get", key))
	case "list":
		value, err := redis.Values(client.Do("lrange", key, 0, 10))
		if err != nil {
			return "", err
		}
		return dealValues(value), nil

	case "set":
		value, err := redis.Values(client.Do("SRANDMEMBER", key, 10))
		if err != nil {
			return "", err
		}
		return dealValues(value), nil

	case "zset":
		value, err := redis.Values(client.Do("zrange", key, 0, 10))
		if err != nil {
			return "", err
		}
		return dealValues(value), nil

	case "hash":
		len, _ := redis.Int(client.Do("hlen", key))
		if len > 10 {
			return "", nil
		}
		value, err := redis.Values(client.Do("hgetall", key))
		if err != nil {
			fmt.Println("hgetall failed", err.Error())
			return "", err
		}
		return dealHashValues(value), nil

	default:
		fmt.Println("unknown key type.")
	}
	return "", nil
}

//DeleteKey get value
func DeleteKey(productName, key string) (int64, error) {
	SetZkInfo(productName)
	client, err := pool.Get()
	if err != nil {
		return 0, err
	}
	defer client.Close()
	return redis.Int64(client.Do("del", key))

}

//Expire expire key
func Expire(productName, key string, time int64) (int64, error) {
	SetZkInfo(productName)
	client, err := pool.Get()
	if err != nil {
		return 0, err
	}
	defer client.Close()
	return redis.Int64(client.Do("expire", key, time))

}

//SetKey get value
func SetKey(productName, key, value string) (string, error) {
	SetZkInfo(productName)
	client, err := pool.Get()
	if err != nil {
		return "", err
	}
	defer client.Close()
	return redis.String(client.Do("set", key, value))

}

//TTL ,get key ttl
func TTL(productName, key string) (int64, error) {
	SetZkInfo(productName)
	client, err := pool.Get()
	if err != nil {
		return 0, err
	}
	defer client.Close()
	return redis.Int64(client.Do("ttl", key))
}

func dealValues(value []interface{}) string {
	fmt.Printf("deal values:%v\n", value)
	var result string
	for _, v := range value {
		sbyte := string(v.([]byte))
		result += sbyte
		result += ","
	}
	return result
}
func dealHashValues(value []interface{}) string {
	fmt.Printf("deal values:%v\n", value)
	var result string
	var i uint16 = 0
	for _, v := range value {
		sbyte := string(v.([]byte))
		result += sbyte

		i++
		if i%2 == 0 {
			result += ","
		} else {
			result += ":"
		}
	}
	return result
}
