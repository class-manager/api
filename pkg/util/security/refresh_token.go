package security

import (
	"crypto/rand"
	"encoding/base32"
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
	db.Conn.Preload("Account").Where(&model.RefreshToken{Token: tokenString}).First(token)

	if token == nil {
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
		Token:      base32.StdEncoding.EncodeToString(tokenStringBytes)[:31],
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
	c.ClearCookie("r_")

	c.Cookie(&fiber.Cookie{
		Name:     "r_",
		Value:    rt,
		Expires:  time.Now().Add(refreshTokenDuration),
		HTTPOnly: true,
		SameSite: "Strict",
		// Domain:   "",
		// Secure:   true,
	})

	return nil
}