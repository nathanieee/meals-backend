package apps

import (
	"fmt"
	"os"
	"os/signal"
	"project-skbackend/configs"
	v1 "project-skbackend/internal/controllers/http/v1"
	"project-skbackend/internal/di"
	"project-skbackend/packages/servers"
	"project-skbackend/packages/utils/utlogger"
	"syscall"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Run(cfg *configs.Config) {
	db, err := gorm.Open(postgres.Open(cfg.DB.GetDbConnectionUrl()), &gorm.Config{
		Logger: logger.Default.LogMode(cfg.GetLogLevel()),
	})

	if err != nil {
		utlogger.LogError(err)
		fmt.Errorf("app - Run - postgres: %w", err)
	}

	err = cfg.DB.DBSetup(db)
	if err != nil {
		utlogger.LogError(err)
		fmt.Errorf("app - Run - DB setup: %w", err)
	}

	rdb := cfg.Redis.GetRedisClient()

	di := di.NewDependencyInjection(db, cfg, rdb)
	handler := gin.New()
	v1.NewRouter(handler, db, cfg, di)
	httpserver := servers.NewServer(handler, servers.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		fmt.Printf("app run: %s", s.String())
	case err := <-httpserver.Notify():
		utlogger.LogError(err)
		fmt.Errorf("%w", err)
	}

	err = httpserver.Shutdown()
	if err != nil {
		utlogger.LogError(err)
		fmt.Errorf("%w", err)
	}
}
