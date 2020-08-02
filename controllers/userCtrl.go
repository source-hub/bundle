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
	bcrypt "golang.org/x/crypto/bcrypt"
	//"bundle/config"
	"bundle/functions"
	"github.com/vladimiroff/jwt-go"

)
type UserController struct{
	Db *gorm.DB
}

func(c UserController)Create(ctx iris.Context){
	var (
		readuser models.User
		checkuser models.User
	)
	err:=ctx.ReadJSON(&readuser)
	if err!=nil{
		ctx.JSON(
			iris.Map{
				"success":false,
				"error":iris.Map{
					"message":"There is a problem creating the user",
				},
		})
		return
	}
	c.Db.Where("email = ?", readuser.Email).First(&checkuser)
	if checkuser.Email != "" {
		ctx.JSON(iris.Map{"success": false, "error": iris.Map{"message": "Username or email already exists"}})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(readuser.Password), 9)
	if err != nil {
		println(err.Error())
		ctx.JSON(iris.Map{
			"success": false,
			"error": iris.Map{
				"message": "There is a problem in creating the user",
			},
		})
		return
	}
	// guid := uuid.New().String()
	// registerUser.Guid = guid
	readuser.Password = string(hash)
	tx := c.Db.Begin()
	dberr := tx.Create(&readuser).Error
	if dberr != nil {
		println(dberr.Error())
		ctx.JSON(iris.Map{
			"success": false,
			"error": iris.Map{
				"message": "There is a problem in creating the user",
			},
		})
		tx.Rollback()
		return
	}
	ctx.JSON(iris.Map{"success": true, "message": "New user created", "user": iris.Map{
		"id":         readuser.ID,
		"email":      readuser.Email,
		"created_at": readuser.CreatedAt,
	}})
	tx.Commit()	
	return	
}
func(c UserController)Login(ctx iris.Context){
	var (
		readuser models.User
		checkuser models.User
	)
	err:=ctx.ReadJSON(&readuser)
	if err!=nil{
		ctx.JSON(iris.Map{
				"success":false,
				"error":iris.Map{
					"message":"body sent is invalid",
				},
			})
		return
	}
	c.Db.Where("email = ?",readuser.Email).First(&checkuser)
	if checkuser.Email == ""{
		ctx.JSON(iris.Map{
				"success":false,
				"error":iris.Map{
					"message":"user not found!",
				},
			})
		return
	}
	hash_err := bcrypt.CompareHashAndPassword([]byte(checkuser.Password), []byte(readuser.Password))
	if hash_err != nil {
		println(hash_err.Error())
		ctx.JSON(iris.Map{
			"success": false,
			"error": iris.Map{
				"message":     "Invalid password",
				"description": hash_err.Error(),
			},
		})
		return

	}
	claims := jwt.MapClaims{
		"id":checkuser.ID,
		"email":checkuser.Email,
	}

	userFunc := functions.UserFunc{DB: c.Db}

	var tokenResponse functions.FunctionResponse
	userFunc.SignToken(claims, &tokenResponse)

	if !tokenResponse.Success {

		ctx.JSON(iris.Map{
			"success": false,
			"error": iris.Map{
				"message":     "There is an issue in logging the user. ",
				"description": tokenResponse.Message,
			},
		})
		return

	}
	ctx.JSON(iris.Map{
		"success": true,
		"token":   tokenResponse.Message,
		"user": iris.Map{
			"email": checkuser.Email,
		},
	})
	return

}
