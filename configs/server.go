package configs

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

func ServerConfig() string {
	serverHost := viper.GetString("SERVER_HOST")
	serverPort := viper.GetString("SERVER_PORT")

	appServer := fmt.Sprintf("%s:%s", serverHost, serverPort)
	log.Print("Server Running at :", appServer)
	return appServer
}
