package apps

import (
	"context"
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
		utlogger.Error(err)
	}

	err = cfg.DB.DBSetup(db)
	if err != nil {
		utlogger.Error(err)
	}

	// * setup context
	ctx := context.Background()

	// * setup redis client
	rdb := cfg.Redis.GetRedisClient()

	// * setup rabbit mq
	ch, close := cfg.Queue.Init()
	defer close()
	cfg.Queue.SetupRabbitMQ(ch, cfg)

	di := di.NewDependencyInjection(db, ch, cfg, rdb, ctx)

	// * setup consumer
	di.ConsumerService.ConsumeTask()

	var forever chan struct{}

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, db, cfg, di)
	server := servers.NewServer(handler, servers.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		utlogger.Info("app run: " + s.String())
	case err := <-server.Notify():
		utlogger.Error(fmt.Errorf("%w", err))
	}

	err = server.Shutdown()
	if err != nil {
		utlogger.Error(fmt.Errorf("%w", err))
	}

	<-forever
}
