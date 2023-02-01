package utils

import "github.com/spf13/viper"

type Config struct {
	DBDriver      string `mapstructure:"DB_DRIVER"`
	DBSource      string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// overwrite the variables got from the config file with the existing environment variables
	viper.AutomaticEnv()
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	viper.Unmarshal(&config)
	return
}
