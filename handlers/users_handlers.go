package handlers

import (
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/espher/GoLang-API-REST/models"
)

type UsersHandler struct {
	db *gorm.DB
}

func UsersRouter(db *gorm.DB) *UsersHandler {
	return &UsersHandler{db: db}
}

func (uh *UsersHandler) GetUsers(c *gin.Context) {
	var users []models.User
	result := uh.db.Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (uh *UsersHandler) GetUserById(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	result := uh.db.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uh *UsersHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if !govalidator.IsEmail(user.Email) {
		c.JSON(http.StatusConflict, gin.H{"error": "Wrong email"})
		return
	}

	// validate email as a unique
	existingUser, err := uh.GetUserByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	//hasnigh psw
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Error  hashing psw"})
		return
	}

	user.Password = string(hash)
	//fmt.Println("user data " + user.Password)
	//create user
	result := uh.db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)

}

func (uh *UsersHandler) UpdateUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	result := uh.db.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if !govalidator.IsEmail(user.Email) {
		c.JSON(http.StatusConflict, gin.H{"error": "Wrong email"})
		return
	}

	result = uh.db.Save(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uh *UsersHandler) DeleteUser(c *gin.Context) {
	// Verifica si el usuario existe en la base de datos
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	result := uh.db.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Intenta eliminar el usuario de la base de datos
	result = result.Delete(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (uh *UsersHandler) GetUserByEmail(email string) (*models.User, error) {
	var existingUser models.User
	result := uh.db.Where("email = ?", email).First(&existingUser)
	if result.Error == nil && result.RowsAffected > 0 {
		return &existingUser, nil
	}
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	return nil, nil
}

func (uh *UsersHandler) GetUserByIdWithoutContext(ID int) (*models.User, error) {
	var existingUser models.User
	result := uh.db.Where("id = ?", ID).First(&existingUser)
	if result.Error == nil && result.RowsAffected > 0 {
		return &existingUser, nil
	}
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	return nil, nil
}
