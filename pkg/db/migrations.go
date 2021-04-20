package db

import (
	model "github.com/class-manager/api/pkg/db/models"
)

// TODO: Implement
func migrate() {
	result := DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if result.Error != nil {
		panic(result.Error)
	}

	DB.AutoMigrate(&model.Account{})
	DB.AutoMigrate(&model.Class{})
	DB.AutoMigrate(&model.Student{})
	DB.AutoMigrate(&model.Lesson{})
	DB.AutoMigrate(&model.BehaviourNote{})
	DB.AutoMigrate(&model.Task{})
	DB.AutoMigrate(&model.TaskResult{})
}
