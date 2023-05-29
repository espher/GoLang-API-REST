package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/espher/GoLang-API-REST/db"
	"github.com/espher/GoLang-API-REST/handlers"
	"github.com/espher/GoLang-API-REST/middlewares"
)

func SetupRoutes(router *gin.Engine) {
	// Users Routes
	usersRouter := handlers.UsersRouter(db.DB)
	router.GET("/users", usersRouter.GetUsers)
	router.GET("/users/:id", usersRouter.GetUserById)
	router.POST("/users", usersRouter.CreateUser)
	router.PUT("/users/:id", usersRouter.UpdateUser)
	router.DELETE("/users/:id", usersRouter.DeleteUser)

	// Posts Routes
	postsRouter := handlers.PostsRouter(db.DB)
	router.GET("/posts", postsRouter.GetPosts)
	router.GET("/posts/:id", postsRouter.GetPostById)
	router.POST("/posts", postsRouter.CreatePost)
	router.PUT("/posts/:id", postsRouter.UpdatePost)
	router.DELETE("/posts/:id", postsRouter.DeletePost)

	loginHandler := handlers.LoginHandlerRoueter(usersRouter)
	router.POST("/login", loginHandler.LoginUser)
	router.POST("/logout", loginHandler.LogoutUser)
	router.GET("/check-area", middlewares.RequireAuth, loginHandler.CheckLogin)

	//General Routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ping",
		})
	})
}
