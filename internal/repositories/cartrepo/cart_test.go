package cartrepo

import (
	"project-skbackend/configs"
	"project-skbackend/mocks"
	"project-skbackend/packages/utils/utlogger"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	gdb *gorm.DB

	cartrepomocks = new(mocks.ICartRepository)
)

func TestMain(m *testing.M) {
	cfg := configs.GetInstance()

	db, err := gorm.Open(postgres.Open(cfg.DB.GetDbConnectionUrl()), &gorm.Config{
		Logger: logger.Default.LogMode(cfg.GetLogLevel()),
	})

	if err != nil {
		utlogger.Error(err)
		panic(err)
	}

	gdb = db

	m.Run()
}
