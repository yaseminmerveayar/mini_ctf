package migrations

import (
	"github.com/yaseminmerveayar/mini_ctf/database"
	"github.com/yaseminmerveayar/mini_ctf/models"
)

func AutoMigrate() error {
	err := database.Database.Db.AutoMigrate(
		&models.User{},
		&models.Question{},
		&models.Category{},
		&models.Log{},
	)
	return err
}
