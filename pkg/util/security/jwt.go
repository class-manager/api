package security

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/pascaldekloe/jwt"
	"github.com/spf13/viper"
)

// CreateJWT creates a new JWT. It requires a UID as well as any other custom
// claims to be added.
//
// If validFor is 0, the JWT will not expire.
func CreateJWT(uid uuid.UUID, validFor time.Duration, customClaims map[string]interface{}) []byte {
	var claims jwt.Claims

	mergedClaims := map[string]interface{}{"uid": uid.String()}
	for k, v := range customClaims {
		mergedClaims[k] = v
	}

	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.Set = mergedClaims
	if validFor > 0 {
		claims.Expires = jwt.NewNumericTime(time.Now().Add(validFor))
	}

	token, _ := claims.HMACSign(jwt.HS512, []byte(viper.GetString("JWT_KEY")))

	return token
}

// ValidateJWT validates a JWT and returns any associated claims.
func ValidateJWT(token []byte) (bool, map[string]interface{}) {
	claims, err := jwt.HMACCheck(token, []byte(viper.GetString("JWT_KEY")))
	if err != nil {
		return false, nil
	}

	// Token has expired
	if !claims.Valid(time.Now()) {
		return false, nil
	}

	return true, claims.Set
}
