package config

import (
	"log"
	"os"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var Cache *redis.Client

func LoadEnvConfigs() {

	log.Println(os.Getenv("APP_ENV"))
	if os.Getenv("APP_ENV") == "production" {
		viper.SetConfigName("prod")
	} else {
		viper.SetConfigName("local")
	}
	viper.SetConfigType("env")
	viper.AddConfigPath("infra/config")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

}

func LoadCache() {
	Cache = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS_HOST"),
		Password: viper.GetString("REDIS_PASSWORD"),
		DB:       viper.GetInt("REDIS_DB"),
	})
}
