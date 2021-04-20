package db

import (
	model "github.com/class-manager/api/pkg/db/models"
)

// TODO: Implement
func migrate() {
	result := Conn.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	if result.Error != nil {
		panic(result.Error)
	}

	Conn.AutoMigrate(&model.Account{})
	Conn.AutoMigrate(&model.Class{})
	Conn.AutoMigrate(&model.Student{})
	Conn.AutoMigrate(&model.Lesson{})
	Conn.AutoMigrate(&model.BehaviourNote{})
	Conn.AutoMigrate(&model.Task{})
	Conn.AutoMigrate(&model.TaskResult{})
	Conn.AutoMigrate(&model.RefreshToken{})
}
