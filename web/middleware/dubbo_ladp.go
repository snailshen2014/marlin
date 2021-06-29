package middleware

import (
	"context"
	"fmt"

	hessian "github.com/apache/dubbo-go-hessian2"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	"github.com/apache/dubbo-go/config"
	"github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/registry/protocol"

	_ "github.com/apache/dubbo-go/filter/filter_impl"

	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

//DubboLdap dubbo service for ldap
type DubboLdap struct {
}

//NewDubboLdap new dubbo ldap service
func NewDubboLdap(cconf config.ConsumerConfig) *DubboLdap {
	config.SetConsumerConfig(cconf)
	dubbo.SetClientConf(dubbo.GetDefaultClientConfig())
	hessian.RegisterPOJO(&LdapUser{})
	config.Load()
	return &DubboLdap{}
}

//NewDubboLdap2 new dubbo ldap service
func NewDubboLdap2() *DubboLdap {
	hessian.RegisterPOJO(&LdapUser{})
	config.Load()
	return &DubboLdap{}
}

// export CONF_PROVIDER_FILE_PATH="xxx"
// export APP_LOG_CONF_FILE="xxx"
// var CONF_CONSUMER_FILE_PATH string = "consumer.yml"
// var APP_LOG_CONF_FILE string = "log.yml"

//CheckUser check user
func (d *DubboLdap) CheckUser(mail, passwd string) (bool, error) {
	args := []interface{}{mail, passwd}
	var checkResult bool
	println("########")
	println(ldapProvider)
	println("########")
	err := ldapProvider.verifyLdapUser(context.TODO(), args, &checkResult)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	println("Response result:%v\n", checkResult)
	return checkResult, nil
}

//GetUser get user
func (d *DubboLdap) GetUser(mail string) (LdapUser, error) {
	args := []interface{}{mail}
	ldapUser := &LdapUser{}
	println("###########")
	println(ldapProvider)
	err := ldapProvider.getUserByName(context.TODO(), args, ldapUser)
	if err != nil {
		fmt.Println(err)
		return *ldapUser, err
	}
	println("Response result:%v\n", ldapUser.name)
	return *ldapUser, nil
}
