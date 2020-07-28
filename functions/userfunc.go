package functions

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/vladimiroff/jwt-go"
	"bundle/models"
	"bundle/config"
)

type UserFunc struct {
	DB *gorm.DB
}

type Claims struct {
	Username string `json:"username"`
	ID       uint   `json:"id"`
	jwt.StandardClaims
}

type FunctionResponse struct {
	Success   bool
	JWTClaims Claims
	User      models.User
	Message   string
}

func (cc *UserFunc) DecryptToken(token string, response *FunctionResponse) bool {

	//
	var claims *Claims

	appConfig := config.Config{}

	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(appConfig.GetConfig()["jwtKey"]), nil
	})

	fmt.Printf("%+v\n", tokenClaims)


	claims, _ = tokenClaims.Claims.(*Claims)

	println("The claims in the jwt are")
	fmt.Printf("%+v\n", claims)
	println("\n\n\n\n")

	if err != nil {
		println("There is an issue with parsing the token")
		println(err.Error())

		response.Success = false
		response.Message = err.Error()

	} else {
		//
		response.Success = true
		response.JWTClaims = *claims
	}

	return true

}

func (cc *UserFunc) SignToken(claims jwt.MapClaims, response *FunctionResponse) bool {

	//

	appConfig := config.Config{}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(appConfig.GetConfig()["jwtKey"]))

	if err != nil {

		response.Success = false
		response.Message = err.Error()
		return false

	}
	response.Success = true
	response.Message = tokenString
	return true

}

func (cc *UserFunc) VerifyToken(token string, response *FunctionResponse) bool {

	//

	cc.DecryptToken(token, response)

	if !response.Success {
		return false
	}
	err := cc.DB.Where("id = ?", response.JWTClaims.ID).First(&response.User).Error
	if err != nil {
		println("There is an issue with getting the user data")
		println(err.Error())
		response.Message = err.Error()
		response.Success = false
	}

	return true
}
