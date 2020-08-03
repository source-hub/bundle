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
		return
	}
	catalog.Root=tokenResponse.User.ID
	catalog.Parent=tokenResponse.User.ID
	catalog.User_id=tokenResponse.User.ID
	dberr:=c.Db.Create(&catalog).Error
	if dberr!=nil{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"problem creating the catalog",
			},
		})
		return
	}
	ctx.JSON(iris.Map{
		"success":true,
		"message":"Catalog created successfully",
		})
	return
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
		return
	}
	var(
		qs models.Query_struct
		ids []uint
		id_e uint
		query_res []models.Query_struct
	)
	cat_type:=ctx.URLParam("type")
	rows,err:=c.Db.Raw(`select id from catalogs where type=? and user_id=?`,cat_type,tokenResponse.User.ID).Rows()
	if err!=nil{
		fmt.Println("the problem is:",err.Error())
		return
	}
	defer rows.Close()
	for rows.Next(){
		err=rows.Scan(&id_e)
		if err!=nil{
			fmt.Println("the problem is:",err.Error())
			return
		}
		ids=append(ids,id_e)
	}
	if len(ids)==0{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"No catalogs found",
			},
		})
		return
	}
	query:=`select 
			catalogs.name as catalog_name,
			catalogs.type as catalog_type,
			item_catalogs.catalog_id as catalog_id,
			items.img as item_img,
			items.id as item_id,
			items.link as item_link 
			from items 
			inner join item_catalogs 
			on items.id=item_catalogs.item_id
			inner join catalogs
			on catalogs.id = catalog_id 
			where catalog_id in (?)`
	rows,err=c.Db.Raw(query,ids).Rows()
	if err!=nil{
		fmt.Println(err,err.Error())
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"problem fetching items",
			},
		})
		return
	}
	defer rows.Close()
	for rows.Next(){
		c.Db.ScanRows(rows,&qs)
		query_res=append(query_res,qs)
	}
	final_response:=make(map[string]models.Custom_Catalog)
	for _,v:=range query_res{
		if cat,ok:=final_response[v.Catalog_name];ok{
			items:=cat.Items
			c_item:=models.Custom_Item{Link:v.Item_link,Img:v.Item_img}
			items=append(items,c_item)
			cat.Items=items
			final_response[v.Catalog_name]=cat
		}else{
			c_cat:=models.Custom_Catalog{Name:v.Catalog_name,Id:v.Catalog_id}
			c_item:=models.Custom_Item{Link:v.Item_link}
			c_cat.Items=[]models.Custom_Item{c_item}
			final_response[v.Catalog_name]=c_cat
		}
	}
	fmt.Println("query_result",final_response)
	ctx.JSON(iris.Map{
		"success":true,
		"message":"Catalogs Fetched successfully",
		"result":final_response,
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
		checkitem models.Item
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
	c.Db.Where("link=?",item.Link).Find(&checkitem)
	if checkitem.ID==0{
		err=c.Db.Create(&item).Error
		if err!=nil{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"prolem creating a record in items table",
			},
		})
		return
		}	
	}else{
		item=checkitem
	}
	crerr:=c.Db.Create(&models.Item_Catalog{Catalog_id:uint(catalog_id),Item_id:item.ID}).Error
	if crerr!=nil{
		fmt.Println(crerr.Error())
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"User already has that item in that catalog already",
			},
		})
		return
	}
	ctx.JSON(iris.Map{
		"success":true,
		"message":"item got added successfully",
	})
}
func (c *CatalogCtrl)To(ctx iris.Context){
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
		checkcat models.Catalog
		readcat models.Catalog
	)
	err:=ctx.ReadJSON(&readcat)
	if err!=nil{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"problem reading the body",
			},
		})
		return
	}
	c.Db.Where("id=?",readcat.ID).Find(&checkcat)
	if checkcat.Name==""{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"catalog does not exist",
			},
		})
		return
	}
	fmt.Println("the object to update",checkcat)
	err=c.Db.Save(&checkcat).Error
	if err!=nil{
		ctx.JSON(iris.Map{
			"success":false,
			"error":iris.Map{
				"message":"error upadting the catalog",
			},
		})
		return
	}
	mess:=fmt.Sprintf("Your catalog is now %s!",readcat.Type)
	ctx.JSON(iris.Map{
		"success":true,
		"message":mess,
	})
}
func (c *CatalogCtrl)Explore(ctx iris.Context){
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
		catalogs []models.Catalog
	)
	query:=`
		select 
		catalogs.name as catalog_name,
		catalogs.type as catalog_type,
		item_catalogs.catalog_id as catalog_id,
		items.img as item_img,
		items.id as item_id,
		items.link as item_link 
		from items 
		inner join item_catalogs 
		on items.id=item_catalogs.item_id
		inner join catalogs
		on catalogs.id = item_catalogs.catalog_id 
		where catalogs.type='public'`
	//resume here
}