package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type LoginHandler struct {
	usersHandler *UsersHandler
}

func LoginHandlerRoueter(usersHandler *UsersHandler) *LoginHandler {
	return &LoginHandler{
		usersHandler: usersHandler,
		// use another dependencies here
	}
}

func (lh *LoginHandler) LoginUser(c *gin.Context) {
	type LoginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// varify content-type
	contentType := c.Request.Header.Get("Content-Type")
	if contentType != "application/json" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content type"})
		return
	}

	// read body request
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}

	// check if is JSON valid
	var loginData LoginData
	if err = json.Unmarshal(body, &loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if !govalidator.IsEmail(loginData.Email) {
		c.JSON(http.StatusConflict, gin.H{"error": "Invalid email"})
		return
	}

	// auth login  logic
	existingUser, err := lh.usersHandler.GetUserByEmail(loginData.Email)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Invalid user or password"})
		return
	}

	if existingUser == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Invalid user or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(loginData.Password))
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Invalid compared hashed psw"})
		return
	}

	userID := int(existingUser.ID)
	fmt.Println(userID)
	fmt.Println(loginData.Password)

	tokenString, err := createToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	//saveo token on cookie
	c.SetSameSite(http.SameSiteLaxMode)
	//false for localhost
	c.SetCookie("Auth", tokenString, 3600*24*15, "", os.Getenv("DOMAIN"), false, true)

	//c.JSON(http.StatusOK, gin.H{"token": tokenString})
	c.JSON(http.StatusOK, gin.H{"token": "Logedin"})
}

// create token with JWT
func createToken(userID int) (string, error) {

	err := godotenv.Load()
	if err != nil {
		panic("Error .env file")
	}

	accessTokenExpireTime := time.Now().Add(time.Hour * 24 * 15).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = accessTokenExpireTime

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (lh *LoginHandler) LogoutUser(c *gin.Context) {
	c.SetCookie("Auth", "", -1, "", os.Getenv("DOMAIN"), false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "logout",
	})
}

func (lh *LoginHandler) CheckLogin(c *gin.Context) {
	var user, _ = c.Get("user")

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
