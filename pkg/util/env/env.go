package env

import "github.com/spf13/viper"

func LoadEnv() {
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "")
	viper.SetDefault("DB_NAME", "classman_prod")
	viper.SetDefault("DB_PORT", 5432)

	viper.AllowEmptyEnv(false)
	viper.AutomaticEnv()
}
