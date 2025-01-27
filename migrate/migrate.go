package main

import (
	"go-crud/initialize"
	"go-crud/models"
)

func init() {
	initialize.LoadEnvVar()
	initialize.ConnectDB()
}

func main() {
	initialize.DB.AutoMigrate(&models.Post{})
}
