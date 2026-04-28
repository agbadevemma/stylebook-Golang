package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/config"
	"github.com/stylebookgo/agbadevemma/internal/db"
)

func main() {
	config.LoadEnv()
	database := db.Connect()
	db.Migrate(database)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	r.Run(port) // listen and serve on
	fmt.Println("Starting server on port", port, "...")
}
