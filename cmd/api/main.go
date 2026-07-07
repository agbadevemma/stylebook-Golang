package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stylebookgo/agbadevemma/internal/config"
	"github.com/stylebookgo/agbadevemma/internal/db"
	"github.com/stylebookgo/agbadevemma/internal/routes"
)

func main() {
	config.LoadEnv()
	database := db.Connect()

	if os.Getenv("RUN_MIGRATIONS") == "true" {
		db.Migrate(database)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000","https://www.stylebookonline.com","https://stylebookonline.com"},
		AllowMethods: []string{
			"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		AllowCredentials: true,
		ExposeHeaders: []string{
			"Content-Length",
		},
		MaxAge: 12 * time.Hour,
	}))

	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Stylebook API is running",
			"time":    time.Now(),
		})
	})
	routes.SetupRoutes(r, database.DB)
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	r.Run(port) // listen and serve on
	fmt.Println("Starting server on port", port, "...")
}
