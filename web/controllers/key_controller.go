/*
 * @Author: your name
 * @Date: 2020-03-08 12:56:58
 * @LastEditTime: 2020-03-12 18:31:25
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/web/controllers/key_controller.go
 */
// file: controllers/user_controller.go

package controllers

import (
	"bufio"
	"fmt"
	"io"
	"marlin/datamodels"
	"marlin/services"
	"strconv"

	"github.com/kataras/iris/v12"
)

// KeyController is our /user controller.
// UserController is responsible to handle the following requests:
// GET  			/user/register
// POST 			/user/register
// GET 				/user/login
// POST 			/user/login
// GET 				/user/me
// All HTTP Methods /user/logout
type KeyController struct {
	// context is auto-binded by Iris on each request,
	// remember that on each incoming request iris creates a new UserController each time,
	// so all fields are request-scoped by-default, only dependency injection is able to set
	// custom fields like the Service which is the same for all requests (static binding)
	// and the Session which depends on the current context (dynamic binding).
	Ctx iris.Context

	// Our KeyService, it's an interface which
	// is binded from the main application.
	Service            services.KeyService
	BatchInfoService   services.UserBatchInfoService
	BatchDetailService services.UserBatchDetailService
	UserService        services.UserService
}

// PostOne handles GET: http://localhost:8080/key/one.
func (c *KeyController) PostOne() datamodels.DbResult {

	var m map[string]interface{}
	// var m KeyJson
	err := c.Ctx.ReadJSON(&m)
	if err != nil {
		fmt.Println("ReadJSON Error:", err)
	}
	fmt.Println(m)
	key, err := c.Service.GetByID(m["key"].(string), m["clusterName"].(string))
	fmt.Printf("PostOne controller ,key:[%v],err:[%v]\n", key, err)
	var (
		msg  string
		code int = 0
	)
	if err != nil {
		if key.Value == "zk" {
			code = -1
			msg = "get key err " + err.Error()
		} else {
			if key.TTL == 0 {
				code = -2
				msg = "key no exists."
			} else {
				code = -1
				msg = "get key err " + err.Error()
			}
		}
	} else {
		if key.Value == "none" {
			code = -2
			msg = "key no exists."
		} else {
			code = 0
			msg = "get key success."
		}

	}
	return *datamodels.NewDbResult(key, code, msg)
	// return id + "sdfsfs111"
}

//ReadFile read upload file's content
func ReadFile(file io.Reader) []string {
	records := make([]string, 100)
	br := bufio.NewReader(file)
	index := 0
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		fmt.Println(string(a))
		records[index] = string(a)
		index++
	}
	return records

}

//PostUpload for upload key's file
func (c *KeyController) PostUpload() datamodels.DbResult {
	// Get the file from the request.
	// file, info, err := c.Ctx.FormFile("uploadfile")
	file, info, err := c.Ctx.FormFile("file")
	if err != nil {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		c.Ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		return *datamodels.NewDbResult("", iris.StatusInternalServerError, "Error while uploading "+err.Error())
	}

	defer file.Close()
	//为了调试，暂时去掉登陆验证
	// user := c.UserService.GetByID(c.getCurrentUserID())
	// if user.Data.(datamodels.User).ID == 0 {
	// 	// if the  session exists but for some reason the user doesn't exist in the "database"
	// 	// then logout and re-execute the function, it will redirect the client to the
	// 	// /user/login page.
	// 	return *datamodels.NewDbResult("", -2, "ok")
	// }
	// userName := user.Data.(datamodels.User).UserName
	fname := info.Filename
	userName := "shenyanjun"
	fmt.Printf("file name:%s\n", fname)
	//judge status
	beforeBatch := c.BatchInfoService.GetByUserAndFileName(userName, fname)
	status := beforeBatch.Data.(datamodels.UserBatchInfo).Status
	beforeFileName := beforeBatch.Data.(datamodels.UserBatchInfo).FileName
	if beforeFileName != "" && status == 0 {
		return *datamodels.NewDbResult("", -3, "ok")
	}
	batchInfo := c.BatchInfoService.Save(datamodels.UserBatchInfo{
		UserName: userName,
		FileName: fname,
		Status:   0,
	})
	batchID := batchInfo.Data.(datamodels.UserBatchInfo).ID
	for _, record := range ReadFile(file) {
		if record != "" {
			fmt.Println("record:" + record)
			//detail
			c.BatchDetailService.Save(datamodels.UserBatchDetail{
				BatchID: int64(batchID),
				KeyName: record,
				Status:  0,
			})
		}

	}
	//goroutine del file
	uploadResponses := make([]interface{}, 0)
	uploadResponses = append(uploadResponses, batchInfo.Data.(datamodels.UserBatchInfo))
	return *datamodels.NewDbResult(uploadResponses, 0, "ok")

}

//GetBatch for upload key's file
func (c *KeyController) GetBatch() datamodels.DbResult {

	//为了调试，暂时去掉登陆校验直接查询shenyanjun用户下所有数据
	// Get the file from the request.
	// user := c.UserService.GetByID(c.getCurrentUserID())
	// if user.Data.(datamodels.User).ID == 0 {
	// 	// if the  session exists but for some reason the user doesn't exist in the "database"
	// 	// then logout and re-execute the function, it will redirect the client to the
	// 	// /user/login page.
	// 	return *datamodels.NewDbResult("", -2, "ok")
	// }
	// userName := user.Data.(datamodels.User).UserName

	// batchInfo := c.BatchInfoService.GetByUserName(userName)
	batchInfo := c.BatchInfoService.GetByUserName("shenyanjun")
	//goroutine del file
	return *datamodels.NewDbResult(batchInfo.Data.([]datamodels.UserBatchInfo), 0, "ok")

}

//GetBatchID ,get batch by batchId
func (c *KeyController) GetBatchId() datamodels.DbResult {

	// Get the file from the request.
	// user := c.UserService.GetByID(c.getCurrentUserID())
	// if user.Data.(datamodels.User).ID == 0 {
	// 	// if the  session exists but for some reason the user doesn't exist in the "database"
	// 	// then logout and re-execute the function, it will redirect the client to the
	// 	// /user/login page.
	// 	return *datamodels.NewDbResult("", -2, "ok")
	// }
	// userName := user.Data.(datamodels.User).UserName
	type MyType struct {
		BatchID int64 `url:"batchID"`
	}
	var query MyType
	err := c.Ctx.ReadQuery(&query)
	if err != nil && !iris.IsErrPath(err) {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		c.Ctx.WriteString(err.Error())
	}
	batchInfo := c.BatchInfoService.GetByBatchId(query.BatchID)

	//goroutine del file
	return *datamodels.NewDbResult(batchInfo.Data.([]datamodels.UserBatchInfo), 0, "ok")

}

//GetDetail for upload key's file
func (c *KeyController) GetDetail() datamodels.DbResult {
	//暂时去掉登陆，为了调试
	// user := c.UserService.GetByID(c.getCurrentUserID())
	// if user.Data.(datamodels.User).ID == 0 {
	// 	// if the  session exists but for some reason the user doesn't exist in the "database"
	// 	// then logout and re-execute the function, it will redirect the client to the
	// 	// /user/login page.
	// 	return *datamodels.NewDbResult("", -2, "ok")
	// }
	type MyType struct {
		BatchID int64 `url:"batchID"`
	}
	var query MyType
	err := c.Ctx.ReadQuery(&query)
	if err != nil && !iris.IsErrPath(err) {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		c.Ctx.WriteString(err.Error())
	}
	fmt.Println(query.BatchID)
	details := c.BatchDetailService.GetDetailstByBatchID(query.BatchID)
	return *datamodels.NewDbResult(details.Data.([]datamodels.UserBatchDetail), 0, "ok")

}

// PostDel handles GET: http://localhost:8080/key/del.
func (c *KeyController) PostDel() datamodels.DbResult {

	var m map[string]interface{}
	// var m KeyJson
	err := c.Ctx.ReadJSON(&m)
	if err != nil {
		fmt.Println("ReadJSON Error:", err)
	}
	fmt.Println(m)
	rtn, err := c.Service.DeleteByID(m["key"].(string), m["clusterName"].(string))
	fmt.Printf("PostDel controller ,rtn:[%v],err:[%v]\n", rtn, err)
	var (
		msg  string
		code int = 0
	)
	if err == nil {
		if rtn == 1 {
			msg = "delete key ok."
		}
		if rtn == 0 {
			code = -2
			msg = "key not exists."
		}
	} else {
		code = -1
		msg = "delete key error." + err.Error()
	}

	return *datamodels.NewDbResult(datamodels.Key{
		ID:    m["key"].(string),
		Value: "",
		TTL:   0,
	}, code, msg)
	// return id + "sdfsfs111"
}

// PostSet handles GET: http://localhost:8080/key/set.
func (c *KeyController) PostSet() datamodels.DbResult {

	var m map[string]interface{}
	// var m KeyJson
	err := c.Ctx.ReadJSON(&m)
	if err != nil {
		fmt.Println("ReadJSON Error:", err)
	}
	fmt.Println(m)
	Key, err := c.Service.Set(datamodels.Key{
		ID:    m["key"].(string),
		Value: m["value"].(string),
	}, m["clusterName"].(string))
	fmt.Printf("PostSet controller ,key:[%v],err:[%v]\n", Key, err)
	var (
		msg  string
		code int = 0
	)
	if err == nil {
		code = 0
	} else {
		code = -1
		msg = "set key error." + err.Error()
	}
	return *datamodels.NewDbResult(datamodels.Key{
		ID:    Key.ID,
		Value: Key.Value,
		TTL:   Key.TTL,
	}, code, msg)
	// return id + "sdfsfs111"
}

//PostExpire ,expire key
func (c *KeyController) PostExpire() datamodels.DbResult {

	var m map[string]interface{}
	// var m KeyJson
	err := c.Ctx.ReadJSON(&m)
	if err != nil {
		fmt.Println("ReadJSON Error:", err)
	}
	fmt.Println(m)
	times := m["time"].(string)
	intTime, _ := strconv.ParseInt(times, 10, 64)
	rtn, err := c.Service.Expire(m["key"].(string), m["clusterName"].(string), intTime)
	fmt.Printf("PostExpire controller ,rtn:[%v],err:[%v]\n", rtn, err)
	var (
		msg  string
		code int = 0
	)
	if err == nil {
		if rtn == 1 {
			msg = "expire key ok."
		}
		if rtn == 0 {
			msg = "key not exists."
			code = -2
		}
	} else {
		code = -1
		msg = "expire key error." + err.Error()
	}

	return *datamodels.NewDbResult(datamodels.Key{
		ID:    m["key"].(string),
		Value: "",
		TTL:   0,
	}, code, msg)
}

//////////////////////////////new
/////////////////////////////////
func (c *KeyController) GetKey() datamodels.DbResult {
	fmt.Printf("GetKey running...,\n")
	type MyType struct {
		Key     string `url:"key"`
		Cluster string `url:"cluster"`
	}
	var query MyType
	err := c.Ctx.ReadQuery(&query)
	if err != nil && !iris.IsErrPath(err) {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		c.Ctx.WriteString(err.Error())
	}
	fmt.Printf("GetKey running...,query=%v\n", query)
	key, err := c.Service.GetByID(query.Key, query.Cluster)
	fmt.Printf("GetKey controller ,key:[%v],err:[%v]\n", key, err)
	var (
		msg  string
		code int = 0
	)
	if err != nil {
		if key.Value == "zk" {
			code = -1
			msg = "get key err " + err.Error()
		} else {
			if key.TTL == 0 {
				code = -2
				msg = "key no exists."
			} else {
				code = -1
				msg = "get key err " + err.Error()
			}
		}
	} else {
		if key.Value == "none" {
			code = -2
			msg = "key no exists."
		} else {
			code = 0
			msg = "get key success."
		}

	}
	//为了vue 使用table
	keys := make([]interface{}, 0)
	keys = append(keys, key)
	return *datamodels.NewDbResult(keys, code, msg)

}

// PostKey handles POST: http://localhost:8080/key/key.
func (c *KeyController) PostKey() datamodels.DbResult {
	fmt.Println("PostKey running...")
	var m map[string]interface{}
	// var m KeyJson
	err := c.Ctx.ReadJSON(&m)
	if err != nil {
		fmt.Println("ReadJSON Error:", err)
	}
	fmt.Println(m)

	Key, err := c.Service.Set(datamodels.Key{
		ID:    m["id"].(string),
		Value: m["value"].(string),
	}, m["clusterName"].(string))
	fmt.Printf("PostKey controller set operator,key:[%v],err:[%v]\n", Key, err)
	if m["ttl"] != nil {
		intTime, _ := strconv.ParseInt(m["ttl"].(string), 10, 64)
		if intTime != -1 {
			// ttlTime := int64(m["ttl"].(float64))
			fmt.Printf("ttlTime:%v\n", intTime)
			rtn, err := c.Service.Expire(m["id"].(string), m["clusterName"].(string), intTime)
			fmt.Printf("PostKey controller expire,rtn:[%v],err:[%v]\n", rtn, err)
		}

	}

	var (
		msg  string
		code int = 0
	)
	if err == nil {
		code = 0
	} else {
		code = -1
		msg = "set key error." + err.Error()
	}
	return *datamodels.NewDbResult(datamodels.Key{
		ID:    Key.ID,
		Value: Key.Value,
		TTL:   Key.TTL,
	}, code, msg)
	// return id + "sdfsfs111"
}

//DeleteKey , delete key
func (c *KeyController) DeleteKey() datamodels.DbResult {
	var m map[string]interface{}
	// var m KeyJson
	err := c.Ctx.ReadJSON(&m)
	if err != nil {
		fmt.Println("ReadJSON Error:", err)
	}
	fmt.Printf("DeleteKey parameters:%v\n", m)
	rtn, err := c.Service.DeleteByID(m["id"].(string), m["clusterName"].(string))
	fmt.Printf("DeleteKey controller ,rtn:[%v],err:[%v]\n", rtn, err)
	var (
		msg  string
		code int = 0
	)
	if err == nil {
		if rtn == 1 {
			msg = "delete key ok."
		}
		if rtn == 0 {
			code = -2
			msg = "key not exists."
		}
	} else {
		code = -1
		msg = "delete key error." + err.Error()
	}

	return *datamodels.NewDbResult(datamodels.Key{
		ID:    m["id"].(string),
		Value: "",
		TTL:   0,
	}, code, msg)
	// return id + "sdfsfs111"
}
