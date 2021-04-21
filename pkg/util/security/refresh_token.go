package security

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/class-manager/api/pkg/db"
	model "github.com/class-manager/api/pkg/db/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

const refreshTokenDuration = time.Hour * 24 * 7

// ValidateRefreshToken validates a refresh token. Returns whether or not the
// refresh token is valid, and an associated uid.
//
// Specify clear=true if the token should be invalidated.
func ValidateRefreshToken(tokenString string, clear bool) (bool, *model.Account) {
	token := new(model.RefreshToken)
	res := db.Conn.Preload("Account").Where(&model.RefreshToken{Token: tokenString}).First(token)

	if res.RowsAffected == 0 {
		return false, nil
	}

	if clear {
		db.Conn.Delete(token)
	}

	return true, &token.Account
}

// createRefreshToken creates a refresh token and stores it in the database.
func createRefreshToken(uid uuid.UUID) (string, error) {
	tokenStringBytes := make([]byte, 25)
	rand.Read(tokenStringBytes)

	token := &model.RefreshToken{
		Token:      fmt.Sprintf("crt_%v", hex.EncodeToString(tokenStringBytes)[:27]),
		AccountID:  uid,
		Expiration: time.Now().Add(refreshTokenDuration),
	}

	res := db.Conn.Create(token)
	if res.Error != nil {
		return "", res.Error
	}

	return token.Token, nil
}

// AddRefreshTokenCookie appends a refresh token to a http response.
func AddRefreshTokenCookie(c *fiber.Ctx, uid uuid.UUID) error {
	rt, err := createRefreshToken(uid)
	if err != nil {
		return err
	}

	// Clear previous cookie
	ClearRefreshCookie(c)

	c.Cookie(&fiber.Cookie{
		Name:     "crt_",
		Value:    rt,
		Expires:  time.Now().Add(refreshTokenDuration),
		HTTPOnly: true,
		SameSite: "None",
		// Domain:   "",
		// Secure:   true,
	})

	return nil
}

// ClearRefreshCookie clears the refresh token cookie.
func ClearRefreshCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "crt_",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HTTPOnly: true,
		SameSite: "None",
		// Domain:   "",
		// Secure:   true,
	})
}
