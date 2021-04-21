package api_v1

import (
	"net/http"

	"github.com/class-manager/api/pkg/db"
	model "github.com/class-manager/api/pkg/db/models"
	"github.com/gofiber/fiber/v2"
)

type registerPostPayload struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func Register(c *fiber.Ctx) error {
	p := new(registerPostPayload)
	if err := c.BodyParser(p); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	// Perform validations
	if err := validate.Struct(p); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	if len(p.Password) < 8 {
		return c.SendStatus(http.StatusBadRequest)
	}

	// Check if a user has already been registered with the email
	acc := &model.Account{Email: p.Email}
	res := db.Conn.Where(acc).First(&model.Account{})
	if res.RowsAffected != 0 {
		return c.SendStatus(http.StatusPreconditionFailed)
	}

	acc = &model.Account{Email: p.Email, Password: calcSHA256([]byte(p.Password)), Name: p.Name}
	res = db.Conn.Create(&acc)
	if res.Error != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	// Account created, authorize
	c.Status(http.StatusCreated)
	return authorize(c, acc)
}
