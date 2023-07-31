package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nathanieiav/project-skbackend/configs"
	routesConfig "github.com/nathanieiav/project-skbackend/routes/configs"
	"github.com/spf13/viper"
)

func main() {
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		viper.AutomaticEnv()
		fmt.Println("Using environment variable")
	} else {
		fmt.Println("Using .env file")
	}

	if err := configs.DBConnection(); err != nil {
		log.Fatalf("database DbConnection error: %s", err)
	}

	router := routesConfig.RouterGroup()
	log.Fatalf("%v", router.Run(configs.ServerConfig()))
}
