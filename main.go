package main

import (
	"rinhaV2/controller"
	"rinhaV2/database"
)

func main() {
	db := database.GetPool()

	controller.InitRoutes(db)
}
