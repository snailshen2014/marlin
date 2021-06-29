/*
 * @Author: david
 * @Date: 2020-03-30 18:42:20
 * @LastEditTime: 2020-11-14 19:33:48
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /marlin/web/middleware/ldapservice.go
 */
package middleware

import (
	"context"
	"time"

	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go/config"
)

var ldapProvider = new(LdapProvider)

func init() {
	config.SetConsumerService(ldapProvider)
	hessian.RegisterPOJO(&LdapUser{})
}

//LdapUser login info
type LdapUser struct {
	id              int32
	name            string
	email           string
	activeDirectory string
	chineseName     string
	phone           string
	expireTime      time.Time
	status          int32
	isDelete        int32
	ts              time.Time
	dn              string
	deparment       string
}

//LdapProvider ,call dubbo service
type LdapProvider struct {
	verifyLdapUser func(ctx context.Context, req []interface{}, rsp *bool) error
	getUserByName  func(ctx context.Context, req []interface{}, rsp *LdapUser) error
}

//Reference reference dubbo service
func (u *LdapProvider) Reference() string {
	return "LdapProvider"
}

//JavaClassName java class name
func (LdapUser) JavaClassName() string {
	return ""

}
