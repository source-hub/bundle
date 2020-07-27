package main


import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"bundle/config"
)

func main(){
	app := iris.New()
	app.WrapRouter(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"POST", "GET", "DELETE"},
		AllowCredentials: true,
	}).ServeHTTP)
	app.Logger().SetLevel("debug")
	config := config.Config{}
	app_config := config.GetConfig()
	customLogger := logger.New(logger.Config{
		Status: true,
		IP: true,
		Method: true,
		Path: true,
		Query: true,

		MessageContextKeys: []string{"logger_message"},
		MessageHeaderKeys: []string{"User-Agent"},
	})

	//db

	db, err := gorm.Open("postgres", app_config["db_env"])

	if err != nil {
		panic(err)
		return 

	}
	db.LogMode(true)
	defer db.Close()

	perform_migrations := true
	drop := false

	if drop {
		Drop(db)
	}

	if perform_migrations {
		Migrate(db)
	}

	//
	app.Use(customLogger)
	api := app.Party("/api/v1")
	api.Get("/health/ping", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "PONG"})
	})
	api.Get("/dev",func(ctx iris.Context){
		ctx.JSON(iris.Map{"message":"lets develop bundle!"})
	})
	port := app_config["PORT"]
	println("The port is:  " + port + "\n\n\n\n")

	app.Run(iris.Addr(":"+port), iris.WithoutServerError(iris.ErrServerClosed))
}