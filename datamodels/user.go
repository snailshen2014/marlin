/*
 * @Author:David
 * @Date: 2020-03-04 21:05:17
 * @LastEditTime: 2020-03-09 00:28:17
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/datamodels/user.go
 */

package datamodels

import "github.com/jinzhu/gorm"

// User is our User example model.
// Keep note that the tags for public-use (for our web app)
// should be kept in other file like "web/viewmodels/user.go"
// which could wrap by embedding the datamodels.User or
// define completely new fields instead but for the shake
// of the example, we will use this datamodel
// as the only one User model in our application.
type User struct {
	gorm.Model
	UserName     string `json:"user_name"`
	AccessKeys   string `json:"access_keys"`
	Organization string `json:"organization"`
	Status       int    `json:"status"`
	RoleType     int    `json:"role_type"`
}
