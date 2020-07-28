package controllers

import (
	//"encoding/json"
	//"fmt"
	//"io/ioutil"
	// "os"
	// "strconv"
	// "strings"
	"bundle/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	iris "github.com/kataras/iris/v12"
	//"bundle/config"
	"bundle/functions"
)

type CatalogCtrl struct{
	Db *gorm.DB
}

func (c *CatalogCtrl)Create(ctx iris.Context){
	token := ctx.GetHeader("token")
	userFunc := functions.UserFunc{DB: c.Db}
	var tokenResponse functions.FunctionResponse
	verifyusr:=userFunc.VerifyToken(token, &tokenResponse)
	println(verifyusr)
	if !verifyusr {
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"could not identify the user",
			},
		})
	}
	var (
		catalog models.Catalog
	)
	err:=ctx.ReadJSON(&catalog)
	if err!=nil{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"problem reading the catalog",
			},
		})
	}
	//resume from here
}