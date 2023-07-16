package configs

import (
	"fmt"

	"github.com/nathanieiav/project-skripsi/models"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConfiguration() string {
	masterDBName := viper.GetString("MASTER_DB_NAME")
	masterDBUser := viper.GetString("MASTER_DB_USER")
	masterDBPassword := viper.GetString("MASTER_DB_PASSWORD")
	masterDBHost := viper.GetString("MASTER_DB_HOST")
	masterDBPort := viper.GetString("MASTER_DB_PORT")
	masterDBSslMode := viper.GetString("MASTER_SSL_MODE")

	masterDBDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		masterDBHost, masterDBUser, masterDBPassword, masterDBName, masterDBPort, masterDBSslMode,
	)

	return masterDBDSN
}

func DBConnection() error {
	var db = DB

	masterDSN := DBConfiguration()
	var err error

	db, err = gorm.Open(postgres.Open(masterDSN), &gorm.Config{})

	if err != nil {
		return err
	}

	err = db.AutoMigrate(migrationModels...)

	DB = db

	return err
}

var migrationModels = []interface{}{
	&models.User{},
}
