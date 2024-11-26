package config

import "github.com/spf13/viper"

type Config struct {
	Port           int    `mapstructure:"PORT"`
	DBPath         string `mapstructure:"DB_PATH"`
	IP             string `mapstructure:"IP"`
	ExternalAPIUrl string `mapstructure:"EXTERNAL_API_URL"`
	MigrationsPath string `mapstructure:"MIGRATIONS_PATH"`
}

func LoadConfig() (cfg *Config, err error) {
	viper.AddConfigPath("../config/envs")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&cfg)

	return
}
