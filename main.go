package main

import (
	"fmt"

	auth_controller "my-project/controller/auth_controller"
	user_controller "my-project/controller/user_controller"
	middleware "my-project/middleware"
	mongo "my-project/mongo"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	var mySlice []int
	mySlice = append(mySlice, 10)
	mySlice = append(mySlice, 20)
	fmt.Println(mySlice)

	router := gin.Default()
	router.Use(CORSMiddleware())

	userDB := mongo.ConnectMongo()

	// auth - service
	router.POST("/register", func(c *gin.Context) {
		auth_controller.Register(c, userDB)
	})
	router.POST("/login", func(c *gin.Context) {
		auth_controller.Login(c, userDB)
	})

	// user - service
	authorized := router.Group("/user", middleware.JWTAuthen())
	authorized.GET("/", func(c *gin.Context) {
		user_controller.GetUsers(c, userDB)
	})
	authorized.GET("/:id", func(c *gin.Context) {
		user_controller.GetUserByID(c, userDB)
	})
	authorized.POST("/", func(c *gin.Context) {
		user_controller.InsertUser(c, userDB)
	})
	authorized.PUT("/:id", func(c *gin.Context) {
		user_controller.UpdateUser(c, userDB)
	})
	authorized.DELETE("/:id", func(c *gin.Context) {
		user_controller.DeleteUser(c, userDB)
	})
	authorized.GET("/profile", func(c *gin.Context) {
		user_controller.GetProfile(c, userDB)
	})

	defer mongo.DisconnectMongo(userDB.DB)
	fmt.Println("===== running in port 8080 =====")
	router.Run(":8080")
}
