package middlewares

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/espher/GoLang-API-REST/db"
	"github.com/espher/GoLang-API-REST/models"
	"github.com/gin-gonic/gin"
)

func RequireAuth(c *gin.Context) {
	fmt.Println("inside middlware")

	//check cookie
	tokenString, err := c.Cookie("Auth")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	fmt.Println(tokenString)

	type MyCustomClaims struct {
		UserID int `json:"user_id"`
		jwt.StandardClaims
	}

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		fmt.Println(err)
		fmt.Println("inside middlware 3")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		//fmt.Println(claims.UserID)
		//fmt.Println(claims.ExpiresAt)

		var user models.User
		db.DB.First(&user, claims.UserID)
		if user.ID == 0 {
			fmt.Println("inside middlware 5")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("user", user)

		c.Next()
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

}
