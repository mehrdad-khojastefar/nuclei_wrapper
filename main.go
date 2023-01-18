package main

import (
	"log"

	"github.com/joho/godotenv"
	"hamravesh.ir/mehrdad-khojastefar/controllers"
	"hamravesh.ir/mehrdad-khojastefar/database"
	"hamravesh.ir/mehrdad-khojastefar/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	// loading global environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err.Error())
	}

	err = database.Db.InitDatabase()
	if err != nil {
		log.Fatal(err.Error())
	}

	r := gin.Default()
	api := r.Group("/api")
	api.POST("/login", controllers.Login)
	api.POST("/register", controllers.Register)

	user := r.Group("/user")
	user.Use(middlewares.JwtAuth())
	user.GET("/", controllers.GetUser)
	user.GET("/jobs", controllers.GetJobs)
	user.POST("/jobs", controllers.AddJob)

	r.Run(":8090")

	// subRunner, err := subfinder.NewRunner("test", &runner.Options{
	// 	Threads:            10,
	// 	Timeout:            30,
	// 	MaxEnumerationTime: 10,
	// 	Resolvers:          resolve.DefaultResolvers,
	// })
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// fmt.Println(subRunner.GetSubdomainArray("iran.ir"))
}
