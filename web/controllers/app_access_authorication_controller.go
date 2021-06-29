/*
 * @Author: your name
 * @Date: 2020-03-09 19:50:06
 * @LastEditTime: 2020-03-09 23:21:46
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /marlin/web/controllers/app_access_authorication_controller.go
 */

package controllers

import (
	"fmt"
	"marlin/datamodels"
	"marlin/services"
	"marlin/utils"

	"github.com/kataras/iris/v12"
)

//AppAccessAuthoricationController application access authorication controller
type AppAccessAuthoricationController struct {
	// context is auto-binded by Iris on each request,
	// remember that on each incoming request iris creates a new UserController each time,
	// so all fields are request-scoped by-default, only dependency injection is able to set
	// custom fields like the Service which is the same for all requests (static binding)
	// and the Session which depends on the current context (dynamic binding).
	Ctx iris.Context
	//Service service for app access
	Service services.AppAccessAuthoricationService
}

//PostRegister ,application access apply
func (c *AppAccessAuthoricationController) PostRegister() datamodels.DbResult {
	//调试去掉登陆校验
	// if !c.isLoggedIn() { //not login
	// 	return loginStaticView
	// }

	var m map[string]interface{}
	// var m KeyJson
	err := c.Ctx.ReadJSON(&m)
	if err != nil {
		fmt.Println("ReadJSON Error:", err)
	}
	fmt.Println(m)

	result := c.Service.Save(datamodels.AppAccessAuthorication{
		ApplicationName: m["applicationName"].(string),
		Applicant:       m["mail"].(string),
		AccessKeys:      m["keyPrefix"].(string),
		Status:          0, //0:未审批，1：已审批，没有审批流程默认通过
		Secret:          "",
	})
	// return mvc.View{
	// 	Name: "menu/application_access_tips.html",
	// 	Data: iris.Map{
	// 		"Application": result.Data.(datamodels.AppAccessAuthorication),
	// 	},
	// }
	return *result
}

//GetRegister RESTful api get application info by application name
func (c *AppAccessAuthoricationController) GetRegister() datamodels.DbResult {
	var m map[string]interface{}
	err := c.Ctx.ReadJSON(&m)
	if err != nil {
		fmt.Println("ReadJSON Error:", err)
	}
	if m["appName"] == "" || m["appName"] == nil {

		return *datamodels.NewDbResult("", -1, "appName parameter error!")
	}

	result := c.Service.GetByName(m["appName"].(string))
	return *result
}

//GetApp ,get my application info
func (c *AppAccessAuthoricationController) GetApp() datamodels.DbResult {
	type QueryParam struct {
		Applicant string `url:"applicant"`
		Status    int64  `url:"status"`
		AppName   string `url:"appName"`
	}
	var query QueryParam
	err := c.Ctx.ReadQuery(&query)
	if err != nil && !iris.IsErrPath(err) {
		c.Ctx.StatusCode(iris.StatusInternalServerError)
		c.Ctx.WriteString(err.Error())
	}

	//query user name,调试暂时去掉用户校验
	// result := c.Service.GetByID(c.getCurrentUserID())
	// user := result.Data.(datamodels.User)
	// if user.ID == 0 {
	// 	return *datamodels.NewDbResult("", -2, "no user")
	// }
	//query user application

	fmt.Println(query)
	//admin
	//需要
	return *c.Service.GetApplications(query.Applicant, query.Status, query.AppName)

}

//PostApproval approval app and return app info
func (c *AppAccessAuthoricationController) PostApproval() datamodels.DbResult {
	var m map[string]interface{}
	// var m KeyJson
	err := c.Ctx.ReadJSON(&m)
	if err != nil {
		fmt.Println("ReadJSON Error:", err)
	}
	fmt.Println(m)

	//query user name
	// result := c.Service.GetByID(c.getCurrentUserID())
	// user := result.Data.(datamodels.User)
	// if user.ID == 0 {
	// 	return *datamodels.NewDbResult("", -2, "no user")
	// }
	// if user.RoleType != 1 {
	// 	return *datamodels.NewDbResult("", -3, "no permission")
	// }
	result2 := c.Service.GetByName(m["appName"].(string))
	app := result2.Data.(datamodels.AppAccessAuthorication)
	app.Secret = utils.GenerateSecret()
	return *c.Service.UpdateByAppName(app)

}
