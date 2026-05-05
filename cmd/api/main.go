package main

import (
	"fmt"
	"os"

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
	routes.SetupRoutes(r, database.DB)
	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	r.Run(port) // listen and serve on
	fmt.Println("Starting server on port", port, "...")
}
