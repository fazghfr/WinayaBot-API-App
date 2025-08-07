package main

import (
	"Discord_API_DB_v1/internal/config"
	"Discord_API_DB_v1/internal/model"
	"github.com/joho/godotenv"
)

func main() {
	// loading environment variables
	err := godotenv.Load()
	if err != nil {
		panic("Env loading failed" + err.Error())
	}

	// init db connection and running migration if needed
	db := config.InitDB()
	err = db.AutoMigrate(&model.User{}, &model.Task{})
	if err != nil {
		panic("Migration failed: " + err.Error())
	}

	// HTTP  server and routes configuration
	Server := config.NewHttpServer()
	Server.SetPort(":8080")
	config.RegisterRoutes(Server)

	Server.Start()
}
