package main


import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
	"log"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"bundle/config"
	"bundle/controllers"
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
	log.Println(app_config["db_env"])
	db, err := gorm.Open("postgres", app_config["db_env"])
	if err != nil {
		log.Fatal("problem connecting to the db") 

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
	log.Println("performed migration!")
	app.Use(customLogger)
	api := app.Party("/api/v1")
	api.Get("/health/ping", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "PONG"})
	})
	//users
	userctrl:=controllers.UserController{Db:db}
	user_apis:=api.Party("/user")
	user_apis.Post("/register",userctrl.Create)
	user_apis.Post("/login",userctrl.Login)

	//catalog

	catalogCtrl:=controllers.CatalogCtrl{Db:db}
	cat_apis:=api.Party("/catalog")
	cat_apis.Post("/create",catalogCtrl.Create)
	cat_apis.Post("/additem",catalogCtrl.AddItem)
	cat_apis.Get("/getcatalogs",catalogCtrl.GetAll)
	cat_apis.Post("/to",catalogCtrl.To)
	cat_apis.Get("/explore",catalogCtrl.Explore)
	api.Get("/dev",func(ctx iris.Context){
		ctx.JSON(iris.Map{"message":"lets develop bundle!"})
	})
	port := app_config["PORT"]
	println("The port is:  " + port + "\n\n\n\n")

	app.Run(iris.Addr(":"+port), iris.WithoutServerError(iris.ErrServerClosed))
}