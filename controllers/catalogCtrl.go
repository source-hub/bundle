package controllers

import (
	//"encoding/json"
	"fmt"
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
		checkcatalog models.Catalog
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
	c.Db.Where("name=? and user_id=?",catalog.Name,tokenResponse.User.ID).Find(&checkcatalog)
	if checkcatalog.ID!=0{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"catalog already exists",
			},
		})
	}
	catalog.Root=tokenResponse.User.ID
	catalog.Parent=tokenResponse.User.ID
	dberr:=c.Db.Create(&catalog).Error
	if dberr!=nil{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"problem creating the catalog",
			},
		})
	}
	ctx.JSON(iris.Map{
		"success":true,
		"message":"Catalog created successfully",
		})
}
func (c *CatalogCtrl)GetAll(ctx iris.Context){
	token := ctx.GetHeader("token")
	userFunc := functions.UserFunc{DB: c.Db}
	var tokenResponse functions.FunctionResponse
	verifyusr:=userFunc.VerifyToken(token, &tokenResponse)
	if !verifyusr {
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"could not identify the user",
			},
		})
	}
	var(
		user_catalogs []models.Catalog
		items []models.Item
	)
	cat_type:=ctx.URLParamString("type")
	c.Db.Where("type=? and user_id=?",cat_type,tokenResponse.User.ID).Find(&user_catalogs)
	if len(user_catalogs==0){
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"No catalogs found",
			},
		})
	}
	query=`select item_catalogs.catalog_id items.name as name,items.link as link from items inner join on items.id=item_catalogs.item_id where catalog_id=?`
	err:=c.Db.Raw(query).Scan(&items).Error
	if err!=nil{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"problem fetching items",
			},
		})
	}
	ctx.JSON(iris.Map{
		"success":true,
		"error":iris.Map{
			"message":"Catalogs Fetched successfully",
		},
	})
}
func (c *CatalogCtrl)AddItem(ctx iris.Context){
	token := ctx.GetHeader("token")
	userFunc := functions.UserFunc{DB: c.Db}
	var tokenResponse functions.FunctionResponse
	verifyusr:=userFunc.VerifyToken(token, &tokenResponse)
	if !verifyusr {
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"could not identify the user",
			},
		})
	}
	var (
		checkcatalog models.Catalog
		item models.Item
		item_cat models.Item_Catalog
	)
	catalog_id:=ctx.URLParamInt("catalog_id")
	c.Db.Where("id=?",catalog_id).Find(&checkcatalog)
	if checkcatalog.Name==""{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"Catalog does not exist",
			},
		})
	}
	err:=ctx.ReadJSON(&item)
	if err!=nil{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"Problem reading the body",
			},
		})
	}
	crerr:=c.Db.Create(&item_cat{Catalog_id:catalog_id,Item_id:item.ID}).Error
	if crerr!=nil{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"Problem creating a record in item_catalogs",
			},
		})
	}
	ctx.JSON(iris.Map{
		"success":true,
		"message":"item got added successfully",
	})
}
	//resume from here