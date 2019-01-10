package main

import (
	"context"
	"github.com/GymWorkoutApp/gwap-files/database"
	"github.com/GymWorkoutApp/gwap-files/models"
)

func main() {
	db := database.GetDB(context.TODO())
	defer db.Close()
	db.AutoMigrate(&models.File{})
}
