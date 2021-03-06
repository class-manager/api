package api_v1

import (
	"github.com/class-manager/api/pkg/db"
	model "github.com/class-manager/api/pkg/db/models"
	"github.com/gofiber/fiber/v2"
)

type dashboardClass struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
}

type dashboardTask struct {
	ID    string         `json:"id"`
	Name  string         `json:"name"`
	Class dashboardClass `json:"class"`
}

type getDashboardReturnPayload struct {
	Classes []*dashboardClass `json:"classes"`
	Tasks   []*dashboardTask  `json:"tasks"`
}

// GetDashboardInfo returns the Classman dashboard info
func GetDashboardInfo(c *fiber.Ctx) error {
	uid := c.Locals("uid")

	classes := new([]model.Class)
	classQuery := db.Conn.Where("account_id = ?", uid).Order("subject_name ASC").Order("name ASC").Find(classes)

	returnClasses := make([]*dashboardClass, 0)

	for _, class := range *classes {
		returnClasses = append(returnClasses, &dashboardClass{
			ID:      class.ID.String(),
			Name:    class.Name,
			Subject: class.SubjectName,
		})
	}

	tasks := new([]model.Task)
	db.Conn.Model(&model.Task{}).Preload("Class").Where("class_id in (?)", classQuery.Select("id")).Order("due_date ASC").Find(tasks)

	returnTasks := make([]*dashboardTask, 0)

	for _, task := range *tasks {
		returnTasks = append(returnTasks, &dashboardTask{
			ID:   task.ID.String(),
			Name: task.Name,
			Class: dashboardClass{
				ID:      task.Class.ID.String(),
				Name:    task.Class.Name,
				Subject: task.Class.SubjectName,
			},
		})
	}

	return c.JSON(&getDashboardReturnPayload{
		Classes: returnClasses,
		Tasks:   returnTasks,
	})
}
