package apps

import (
	"fmt"
	"os"
	"os/signal"
	"project-skbackend/configs"
	v1 "project-skbackend/internal/controllers/http/v1"
	"project-skbackend/internal/di"
	"project-skbackend/packages/servers"
	"syscall"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Run(cfg *configs.Config) {
	db, err := gorm.Open(postgres.Open(cfg.DB.GetDbConnectionUrl()))
	if err != nil {
		fmt.Errorf("app - Run - postgres: %w", err)
	}

	err = cfg.DB.AutoSeedEnum(db)
	if err != nil {
		fmt.Errorf("app - Run - create enum: %w", err)
	}

	err = cfg.DB.AutoMigrate(db)
	if err != nil {
		fmt.Errorf("app - Run - migrate: %w", err)
	}

	err = cfg.DB.AutoSeedTable(db)
	if err != nil {
		fmt.Errorf("app - Run - seed table: %w", err)
	}

	di := di.NewDependencyInjection(db, cfg)
	handler := gin.New()
	v1.NewRouter(handler, db, cfg, di)
	httpServer := servers.NewServer(handler, servers.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		fmt.Printf("app run: %s", s.String())
	case err := <-httpServer.Notify():
		fmt.Errorf("%w", err)
	}

	err = httpServer.Shutdown()
	if err != nil {
		fmt.Errorf("%w", err)
	}
}
