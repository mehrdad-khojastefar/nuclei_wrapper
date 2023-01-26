package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// loading global environment variables
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	err := Db.InitDatabase()
	if err != nil {
		log.Fatal(err.Error())
	}

	r := gin.Default()
	api := r.Group("/api")
	api.POST("/login", Login)
	api.POST("/register", Register)

	user := r.Group("/user")
	user.Use(JwtAuth())
	user.GET("/", GetUser)
	user.GET("/jobs/*id", GetJob)
	user.GET("/job/:id/*action", ManageJob)
	user.POST("/jobs", AddJob)

	r.Run(":8000")
}
