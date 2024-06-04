package main

import (
	"service-catalog/config"
	"service-catalog/internal/db"
	"service-catalog/internal/router"
)

func main() {
	config.LoadConfig()
	port := config.GetEnv("PORT", "8080")

	client, err := db.SetupDatabase()
	if err != nil {
		panic(err)
	}
	r := router.SetupRouter(client)
	r.Run(":" + port)
}
