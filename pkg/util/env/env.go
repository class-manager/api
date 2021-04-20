package env

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/spf13/viper"
)

func LoadEnv() {
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_NAME", "classman_prod")
	viper.SetDefault("DB_PORT", 5432)

	viper.SetDefault("JWT_KEY", randomHex(32))

	viper.AllowEmptyEnv(false)
	viper.AutomaticEnv()
}

func randomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}
