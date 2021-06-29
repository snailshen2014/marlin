package main

import (
	"marlin/repositories"
	"marlin/services"
	"marlin/tasks"
	"marlin/web/controllers"
	"marlin/web/middleware"

	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

func main() {
	app := newApp()
	mvcHandle(app)
	app.Run(
		// Starts the web server at localhost:8080
		iris.Addr("localhost:8088"),
		// Ignores err server closed log when CTRL/CMD+C pressed.
		iris.WithoutServerError(iris.ErrServerClosed),
		// Enables faster json serialization and more.
		iris.WithOptimizations,
	)
}

func newApp() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel("debug")
	//vue static load
	//register static file
	// app.HandleDir("/static", "./web/static")
	// //register view file
	// app.RegisterView(iris.HTML("./web/", ".html"))
	// app.Handle("GET", "/", func(ctx iris.Context) {
	// 	ctx.View("index.html")
	// })
	// app.Configure(iris.WithConfiguration(iris.Configuration{
	// 	Charset: "UTF-8",
	// }))

	//test jwt
	// app.Get("/getJWT", func(ctx iris.Context) {
	// 	//set jwt
	// 	token := jwt.NewTokenWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 		"userID":   "123abc",
	// 		"userName": "shenyanjun",
	// 		"admin":    true,
	// 		// 签发人
	// 		"iss": "Marlin",
	// 		// 签发时间
	// 		"iat": time.Now().Unix(),
	// 		// 设定过期时间，便于测试，设置1分钟过期
	// 		"exp": time.Now().Add(1 * time.Minute * time.Duration(1)).Unix(),
	// 	})
	// 	//签名生成jwt字符串
	// 	tokenString, _ := token.SignedString([]byte(middleware.GloablJwtSecret))
	// 	ctx.JSON(tokenString)
	// })
	//showHello加上jwt验证
	//在请求头里加入
	//Authorization Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTg5NzIyMzI1LCJpYXQiOjE1ODk3MjIyNjUsImlzcyI6Ik1hcmxpbiIsInVzZXJJRCI6IjEyM2FiYyIsInVzZXJOYW1lIjoic2hlbnlhbmp1biJ9.ZDvw71dbMbZfPsF4jYZHmy6QXWit1FPCfrw6BnnW4NU
	// app.Get("/showHello", middleware.NewJwtAuth().GetJwt().Serve, func(ctx iris.Context) {
	// 	ctx.JSON("Hello Iris JWT")
	// })
	// app.Get("/showJWT", middleware.NewJwtAuth().GetJwt().Serve, func(ctx iris.Context) {
	// 	jwtInfo := ctx.Values().Get("jwt").(*jwt.Token)
	// 	//获取签发jwt写入的数据
	// 	userID := jwtInfo.Claims.(jwt.MapClaims)["userID"].(string)
	// 	fmt.Printf("Jwt userID:%s\n", userID)
	// 	ctx.JSON(jwtInfo)

	// })

	return app
}

func mvcHandle(app *iris.Application) {
	//session, not using session ,jwt replace it
	// sessionFactory := middleware.NewSessionFactory()
	// app.Use(sessionFactory.GetSession().Handler())

	//cors
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, //允许通过的主机名称
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "HEAD", "DELETE", "OPTIONS"},
	})
	// //use cors all router
	app.Use(crs)
	// app.Use(middleware.NewJwtAuth().GetJwt().Serve)
	repo := repositories.NewUserRepository()
	userService := services.NewUserService(repo)

	batchInfoRepo := repositories.NewUserBatchInfoRepository()
	batchInfoService := services.NewUserBatchInfoService(batchInfoRepo)
	batchDetailRepo := repositories.NewUserBatchDetailRepository()
	batchDetailService := services.NewUserBatchDetailService(batchDetailRepo)

	keyRepo := repositories.NewKeyRepository()
	keyService := services.NewKeyService(keyRepo)
	keyMvcApp := mvc.New(app.Party("/key", middleware.NewJwtAuth().GetJwt().Serve).AllowMethods(iris.MethodOptions))
	keyMvcApp.Register(
		keyService,
		batchInfoService,
		batchDetailService,
		userService,
		// sessionFactory.GetSession().Start,
	)
	keyMvcApp.Handle(new(controllers.KeyController))

	appAccessRepo := repositories.NewAppAccessAuthoricationRepository()
	appAccessService := services.NewAppAccessAuthoricationService(appAccessRepo)
	applicationAccess := mvc.New(app.Party("/application", middleware.NewJwtAuth().GetJwt().Serve).AllowMethods(iris.MethodOptions)).Register(appAccessService)
	applicationAccess.Handle(new(controllers.AppAccessAuthoricationController))
	//users no set jwt auth
	userMvcApp := mvc.New(app.Party("/users").AllowMethods(iris.MethodOptions))
	userMvcApp.Register(
		userService,
		// sessionFactory.GetSession().Start,
	)
	userMvcApp.Handle(new(controllers.UserController))

	//deal key
	tasks.DeleteKeys(batchInfoService, batchDetailService, keyService)

}
