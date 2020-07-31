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
	//controller for creating a catalog

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
		return
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
		return
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
	return
}
func (c *CatalogCtrl)GetAll(ctx iris.Context){
	//
	//controller for fetching all catalogs with items
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
		return
	}
	type query_struct struct{
		Catalog_id uint `json:"catalog_id"`
		models.Item
	}
	var(
		ids []uint
		//user_catalogs []models.Catalog
		//items []models.Item
		query_res []query_struct
	)
	cat_type:=ctx.URLParam("type")
	c.Db.Select("id").Where("type=? and user_id=?",cat_type,tokenResponse.User.ID).Find(&ids)
	if len(ids)==0{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"No catalogs found",
			},
		})
	}
	query:=`select 
			item_catalogs.catalog_id,
			items.img as img,
			items.link as link 
			from items 
			inner join item_catalogs 
			on items.id=item_catalogs.item_id 
			where catalog_id=?`
	err:=c.Db.Raw(query).Scan(&query_res).Error
	if err!=nil{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"problem fetching items",
			},
		})
	}
	fmt.Println("query_result",query_res)
	ctx.JSON(iris.Map{
		"success":true,
		"message":"Catalogs Fetched successfully",
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
		return
	}
	var (
		checkcatalog models.Catalog
		item models.Item
	)
	catalog_id,err:=ctx.URLParamInt("catalog_id")
	if err!=nil{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"Please provide catalog_id",
			},
		})
		return
	}
	c.Db.Where("id=?",catalog_id).Find(&checkcatalog)
	if checkcatalog.Name==""{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"Catalog does not exist",
			},
		})
		return
	}
	err=ctx.ReadJSON(&item)
	if err!=nil{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"Problem reading the body",
			},
		})
		return
	}
	err=c.Db.Create(&item).Error
	if err!=nil{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"Problem creating a record in items",
			},
		})
		return
	}
	crerr:=c.Db.Create(&models.Item_Catalog{Catalog_id:uint(catalog_id),Item_id:item.ID}).Error
	if crerr!=nil{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"Problem creating a record in item_catalogs",
			},
		})
		return
	}
	ctx.JSON(iris.Map{
		"success":true,
		"message":"item got added successfully",
	})
}
	//resume from here