package config

import (
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/viper"
)

var cfg *config

type config struct {
	WebServerName     string `mapstructure:"WEBSERVER_NAME"`
	WebServerPort     string `mapstructure:"WEBSERVER_PORT"`
	RedisHost         string `mapstructure:"REDIS_HOST"`
	RedisPort         string `mapstructure:"REDIS_PORT"`
	RedisDB           int    `mapstructure:"REDIS_DB"`
	RedisPassword     string `mapstructure:"REDIS_PASSWORD"`
	LimitRequestIP    int    `mapstructure:"LIMIT_REQUEST_IP"`
	LimitRequestToken int    `mapstructure:"LIMIT_REQUEST_TOKEN"`
	IpBlockingTime    int    `mapstructure:"IP_BLOCKING_TIME"`
	TokenBlockingTime int    `mapstructure:"TOKEN_BLOCKING_TIME"`
}

func NewConfig() *config {
	err := cfg.validate()
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}

func init() {
	viper.AutomaticEnv()

	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.ReadInConfig()

	viper.BindEnv("WEBSERVER_NAME")
	viper.BindEnv("WEBSERVER_PORT")
	viper.BindEnv("REDIS_HOST")
	viper.BindEnv("REDIS_PORT")
	viper.BindEnv("REDIS_DB")
	viper.BindEnv("REDIS_PASSWORD")
	viper.BindEnv("LIMIT_REQUEST_IP")
	viper.BindEnv("LIMIT_REQUEST_TOKEN")
	viper.BindEnv("IP_BLOCKING_TIME")
	viper.BindEnv("TOKEN_BLOCKING_TIME")

	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal(err)
	}

}

func (c *config) validate() error {
	t := reflect.TypeOf(*cfg)
	v := reflect.ValueOf(*cfg)

	for i := range t.NumField() {
		value := v.Field(i).Interface()

		if value == "" {
			return fmt.Errorf("enviroment variable not found: %v", t.Field(i).Name)
		}
	}

	return nil
}
