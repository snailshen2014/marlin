package middleware

import (
	"github.com/iris-contrib/middleware/jwt"
	jailer "github.com/iris-contrib/middleware/jwt"
)

//JwtAuth jwt auth
type JwtAuth struct {
	jwt *jailer.Middleware
}

//GloablJwtSecret gloabl marlin jwt secret
var GloablJwtSecret string = "marlin-secret"

//NewJwtAuth new JwtAuth
func NewJwtAuth() *JwtAuth {
	return &JwtAuth{
		jwt: jwt.New(jwt.Config{
			// Extractor属性可以选择从什么地方获取jwt进行验证，默认从http请求的header中的Authorization字段提取，也可指定为请求参数中的某个字段
			// 从请求参数token中提取
			// Extractor: jwt.FromParameter("token"),
			// 从请求头的Authorization字段中提取，这个是默认值
			Extractor: jwt.FromAuthHeader,
			// 注意，新增了一个错误处理函数
			// ErrorHandler: func(ctx context.Context, err error) {
			// 	if err == nil {
			// 		return
			// 	}

			// 	ctx.StopExecution()
			// 	ctx.StatusCode(iris.StatusUnauthorized)
			// 	ctx.JSON(ResModel{
			// 		Code: "501",
			// 		Msg:  err.Error(),
			// 	})
			// },
			// 设置一个函数返回秘钥，关键在于return []byte("这里设置秘钥")
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return []byte(GloablJwtSecret), nil
			},
			// 设置一个加密方法
			SigningMethod: jwt.SigningMethodHS256,
		}),
	}
}

//GetJwt get sessopm
func (sf *JwtAuth) GetJwt() *jailer.Middleware {
	return sf.jwt
}
