package handlers

import (
	"net/http"
	"strconv"

	"github.com/espher/GoLang-API-REST/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostsHandler struct {
	db *gorm.DB
}

func PostsRouter(db *gorm.DB) *PostsHandler {
	return &PostsHandler{db: db}
}

func (ph *PostsHandler) GetPosts(c *gin.Context) {
	var posts []models.Post
	//result := ph.db.Find(&posts).Joins("join users on posts.user_id = users.id")
	result := ph.db.Preload("User").Find(&posts)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (ph *PostsHandler) GetPostById(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post
	result := ph.db.First(&post, postID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (ph *PostsHandler) CreatePost(c *gin.Context) {
	var postData models.Post
	if err := c.ShouldBindJSON(&postData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post data"})
		return
	}

	// Verificar si el usuario existe
	var user models.User
	result := ph.db.First(&user, postData.UserID)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	// Crear el nuevo post
	newPost := models.Post{
		Title:       postData.Title,
		Description: postData.Description,
		UserID:      postData.UserID,
	}
	result = ph.db.Create(&newPost)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusOK, newPost)
}

func (ph *PostsHandler) UpdatePost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post
	result := ph.db.First(&post, postID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	result = ph.db.Save(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (ph *PostsHandler) DeletePost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post
	result := ph.db.First(&post, postID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	result = ph.db.Delete(&post)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
