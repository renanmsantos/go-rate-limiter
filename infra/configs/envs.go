package configs

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func LoadEnvConfigs() {

	log.Println(os.Getenv("APP_ENV"))
	if os.Getenv("APP_ENV") == "production" {
		viper.SetConfigName("prod")
	} else {
		viper.SetConfigName("local")
	}
	viper.SetConfigType("env")
	viper.AddConfigPath("infra/configs")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}
