package config

import "github.com/spf13/viper"

func NewViper() *viper.Viper {

	config := viper.New()
	config.SetConfigName("dev")
	config.SetConfigType("yaml")
	config.AddConfigPath(".")

	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return config
}
