// file: controllers/user_controller.go

package controllers

import (
	"fmt"
	"marlin/datamodels"
	"marlin/services"
	"marlin/web/middleware"
	"time"

	"github.com/iris-contrib/middleware/jwt"
	"github.com/kataras/iris/v12"
)

// UserController is our /user controller.
// UserController is responsible to handle the following requests:
// GET  			/user/register
// POST 			/user/register
// GET 				/user/login
// POST 			/user/login
// GET 				/user/me
// All HTTP Methods /user/logout
type UserController struct {
	// context is auto-binded by Iris on each request,
	// remember that on each incoming request iris creates a new UserController each time,
	// so all fields are request-scoped by-default, only dependency injection is able to set
	// custom fields like the Service which is the same for all requests (static binding)
	// and the Session which depends on the current context (dynamic binding).
	Ctx iris.Context

	// Our UserService, it's an interface which
	// is binded from the main application.
	Service services.UserService
}

// const userIDKey = "UserID"
// const isLoggedInKey = "IsLoggedIn"

// func (c *UserController) getCurrentUserID() int64 {
// 	session := sessions.Get(c.Ctx)
// 	userID, _ := session.GetInt64(userIDKey)
// 	return userID
// }

// func (c *UserController) setCurrentUserID(userID int64) {
// 	session := sessions.Get(c.Ctx)
// 	session.Set(userIDKey, userID)
// }

// func (c *UserController) isLoggedIn() bool {
// 	session := sessions.Get(c.Ctx)
// 	loggedIn, _ := session.GetBoolean(isLoggedInKey)
// 	return loggedIn
// }

// func (c *UserController) logout() {
// 	session := sessions.Get(c.Ctx)
// 	session.Set(isLoggedInKey, false)
// 	session.Destroy()
// }
// func (c *UserController) login() {
// 	session := sessions.Get(c.Ctx)
// 	session.Set(isLoggedInKey, true)
// }

//GetUser  get user's info
func (c *UserController) GetUser() datamodels.DbResult {
	type QueryParam struct {
		UserName string `url:"name"`
	}
	var query QueryParam
	err := c.Ctx.ReadQuery(&query)
	if err != nil && !iris.IsErrPath(err) {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		c.Ctx.WriteString(err.Error())
	}
	fmt.Println("########" + query.UserName)

	result := c.Service.GetByName(query.UserName)
	return result

}

//GetApproval ,approval and return user info
func (c *UserController) GetApproval() datamodels.DbResult {
	type MyType struct {
		Flag string `url:"flag"`
	}
	var query MyType
	err := c.Ctx.ReadQuery(&query)
	if err != nil && !iris.IsErrPath(err) {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		c.Ctx.WriteString(err.Error())
	}
	fmt.Println("GetApproval ########" + query.Flag)
	//query user name
	result := c.Service.GetByID(1)
	user := result.Data.(datamodels.User)
	if user.ID == 0 {
		return *datamodels.NewDbResult("", -2, "no user")
	}
	//admin
	if user.RoleType != 1 {
		return *datamodels.NewDbResult("", -3, "no permission")
	}
	//query user name
	result2 := c.Service.GetByNameOrStatus(query.Flag, -1)
	return c.Service.UpdateByName(result2.Data.(datamodels.User))

}

//PostLogin user login
func (c *UserController) PostLogin() datamodels.DbResult {
	var m map[string]interface{}
	// var m KeyJson
	err := c.Ctx.ReadJSON(&m)
	if err != nil {
		fmt.Println("ReadJSON Error:", err)
	}
	fmt.Println(m)

	var (
		msg  string
		code int = 0
	)
	result := c.Service.GetByName(m["mail"].(string))
	user := result.Data.(datamodels.User)
	if user.ID == 0 {
		//用户不存在
		code = -1
		msg = "用户不存在,请注册后登录！"
		return *datamodels.NewDbResult(datamodels.User{
			UserName: m["mail"].(string),
			Status:   user.Status,
		}, code, msg)
	}

	//no Approval and then give tips,0:register,1:passed
	if user.Status == 0 {
		//用户还未审核
		code = -2
		msg = "用户还未审核，请审核通过后登录！"
		return *datamodels.NewDbResult(datamodels.User{
			UserName: m["mail"].(string),
			Status:   user.Status,
		}, code, msg)
	}

	if m["password"].(string) != "" {
		// config.ConsumerInit("consumer.yml")
		// dubboLdap := middleware.NewDubboLdap(config.GetConsumerConfig())
		// dubboLdap := middleware.NewDubboLdap2()
		//check password
		// check, _ := dubboLdap.CheckUser(mail, password)
		// if !check {
		// 	fmt.Printf("mail:%s, password err.\n", mail)
		// 	return nil
		// }

		// _, err := dubboLdap.GetUser(mail)
		// println(err)

		//用户密码不对
		// code = -3
		// msg = "用户密码错误！"
	}
	var isAdmin bool
	if user.RoleType == 1 {
		isAdmin = true
	} else {
		isAdmin = false
	}
	jwtString := GetJWTJsonString(int64(user.ID), user.UserName, isAdmin)
	type token struct {
		Token   string `json:"token"`
		IsAdmin bool   `json:"isAdmin"`
	}
	return *datamodels.NewDbResult(token{
		Token:   jwtString,
		IsAdmin: isAdmin,
	}, code, msg)
}

// Post handles POST: http://localhost:8080/users.
func (c *UserController) Post() datamodels.DbResult {
	var m map[string]interface{}
	// var m KeyJson
	err := c.Ctx.ReadJSON(&m)
	if err != nil {
		fmt.Println("ReadJSON Error:", err)
	}
	fmt.Println(m)
	var (
		msg  string
		code int = 0
	)

	result := c.Service.GetByName(m["mail"].(string))
	user := result.Data.(datamodels.User)
	if user.ID != 0 {
		msg = "用户已经存在！"
		code = -1
		return *datamodels.NewDbResult(datamodels.User{
			UserName: m["mail"].(string),
			Status:   user.Status,
		}, code, msg)
	}

	// create the new user, the password will be hashed by the service.
	_, err2 := c.Service.Save(datamodels.User{
		UserName:     m["mail"].(string),
		AccessKeys:   m["keyPrefix"].(string),
		Organization: m["project"].(string),
		Status:       0,
	})
	if err2 != nil {
		msg = "系统错误"
		code = -2
		return *datamodels.NewDbResult(datamodels.User{
			UserName: m["mail"].(string),
			Status:   user.Status,
		}, code, msg)
	}
	msg = "注册成功，请等待审核！"
	return *datamodels.NewDbResult(datamodels.User{
		UserName: m["mail"].(string),
		Status:   user.Status,
	}, code, msg)

}

//Get ,get user by user's name and status
func (c *UserController) Get() datamodels.DbResult {
	type QueryParam struct {
		Name   string `url:"name"`
		Status int64  `url:"status"`
	}
	var query QueryParam
	err := c.Ctx.ReadQuery(&query)
	if err != nil && !iris.IsErrPath(err) {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		c.Ctx.WriteString(err.Error())
	}
	//query user name
	//测试指定user id
	var tempId int64 = 2
	// result := c.Service.GetByID(c.getCurrentUserID())
	result := c.Service.GetByID(tempId)
	user := result.Data.(datamodels.User)
	if user.ID == 0 {
		return *datamodels.NewDbResult("", -2, "no user")
	}
	//admin
	if user.RoleType == 1 {
		return c.Service.GetByNameOrStatus(query.Name, query.Status)
	}
	//only query loggin user
	return c.Service.GetByNameOrStatus(user.UserName, query.Status)

}

//PostApproval, approval user
func (c *UserController) PostApproval() datamodels.DbResult {
	var m map[string]interface{}
	// var m KeyJson
	err := c.Ctx.ReadJSON(&m)
	if err != nil {
		fmt.Println("ReadJSON Error:", err)
	}
	fmt.Println(m)

	//query user name
	// result := c.Service.GetByID(c.getCurrentUserID())
	var tempId int64 = 2
	result := c.Service.GetByID(tempId)
	user := result.Data.(datamodels.User)
	if user.ID == 0 {
		return *datamodels.NewDbResult("", -2, "no user")
	}
	//admin
	if user.RoleType != 1 {
		return *datamodels.NewDbResult("", -3, "no permission")
	}
	//query user name
	fmt.Println(m["ID"].(float64))

	result2 := c.Service.GetByID(int64(m["ID"].(float64)))
	return c.Service.UpdateByName(result2.Data.(datamodels.User))

}

//GetJWTJsonString generate jwt string
func GetJWTJsonString(userID int64, userName string, isAdmin bool) string {
	//set jwt
	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":   userID,
		"userName": userName,
		"admin":    isAdmin,
		// 签发人
		"iss": "Marlin",
		// 签发时间
		"iat": time.Now().Unix(),
		// 设定过期时间，便于测试，设置1天过期
		"exp": time.Now().Add(24 * time.Hour * time.Duration(1)).Unix(),
	})
	//签名生成jwt字符串
	tokenString, _ := token.SignedString([]byte(middleware.GloablJwtSecret))
	return tokenString
}
