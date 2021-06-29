/*
 * @Author: David
 * @Date: 2020-03-04 21:05:17
 * @LastEditTime: 2020-11-14 19:38:11
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /mvc-demo/datamodels/user.go
 */

package datamodels

import "github.com/jinzhu/gorm"

// AppAccessAuthorication is our User example model.
// Keep note that the tags for public-use (for our web app)
// should be kept in other file like "web/viewmodels/user.go"
// which could wrap by embedding the datamodels.User or
// define completely new fields instead but for the shake
// of the example, we will use this datamodel
// as the only one User model in our application.
type AppAccessAuthorication struct {
	gorm.Model
	ApplicationName string `json:"application_name"`
	Secret          string `json:"secret"`
	Applicant       string `json:"applicant"`
	AccessKeys      string `json:"access_keys"`
	Status          int    `json:"status"`
}
