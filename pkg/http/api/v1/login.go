package api_v1

import (
	"crypto/sha256"
	"net/http"
	"time"

	"github.com/class-manager/api/pkg/db"
	model "github.com/class-manager/api/pkg/db/models"
	"github.com/class-manager/api/pkg/util/security"
	"github.com/gofiber/fiber/v2"
)

type loginPostPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}

func Login(c *fiber.Ctx) error {
	p := new(loginPostPayload)

	// Check that payload is valid
	if err := c.BodyParser(p); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	// Perform validations
	if err := validate.Struct(p); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	// We have an email and password
	acc := new(model.Account)
	hash := calcSHA256([]byte(p.Password))
	res := db.Conn.Where(&model.Account{Email: p.Email, Password: hash}).First(acc)

	if res.RowsAffected == 0 {
		return c.SendStatus(http.StatusNotFound)
	}

	// Valid, authorize session
	return authorize(c, acc)
}

func calcSHA256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

type authorizeResponsePayload struct {
	Token string `json:"token"`
}

func authorize(c *fiber.Ctx, account *model.Account) error {
	uid := account.ID

	token := string(security.CreateJWT(uid, time.Minute*15, make(map[string]interface{})))
	security.AddRefreshTokenCookie(c, uid)

	return c.JSON(&authorizeResponsePayload{Token: token})
}
