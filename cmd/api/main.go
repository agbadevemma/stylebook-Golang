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

	// db.Migrate(database)
	// if os.Getenv("RUN_MIGRATIONS") == "true" {
	// }

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
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

	routes.SetupRoutes(r, database.DB)
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	r.Run(port) // listen and serve on
	fmt.Println("Starting server on port", port, "...")
}
