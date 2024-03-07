package main

import (
	"context"

	"github.com/Ferriem/go-web/Product/common"
	"github.com/Ferriem/go-web/Product/fronted/middleware"
	"github.com/Ferriem/go-web/Product/fronted/web/controllers"
	"github.com/Ferriem/go-web/Product/repositories"
	"github.com/Ferriem/go-web/Product/services"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/opentracing/opentracing-go/log"
)

func main() {
	// create iris application
	app := iris.New()
	// set error mode
	app.Logger().SetLevel("debug")
	//register template
	template := iris.HTML("./fronted/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)
	// register template target
	app.HandleDir("/assets", "./fronted/web/assets")
	app.HandleDir("/html", "./fronted/web/htmlProductShow")
	//error handling
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("Message", ctx.Values().GetStringDefault("message", "The page you're looking for doesn't exist"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	//connect to mysql
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//register controllers
	userRepository := repositories.NewUserManager("user", db)
	userService := services.NewUserService(userRepository)
	userParty := app.Party("/user")
	user := mvc.New(userParty)
	user.Register(ctx, userService)
	user.Handle(new(controllers.UserController))

	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderManager("`order`", db)
	orderService := services.NewOrderService(order)
	productParty := app.Party("/product")
	productMvc := mvc.New(productParty)
	productParty.Use(middleware.AuthConnProduct)
	productMvc.Register(ctx, productService, orderService)
	productMvc.Handle(new(controllers.ProductController))
	// start the server
	app.Run(
		iris.Addr("localhost:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
