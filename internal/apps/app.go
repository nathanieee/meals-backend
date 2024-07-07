package apps

import (
	"context"
	"os"
	"os/signal"
	"project-skbackend/configs"
	v1 "project-skbackend/internal/controllers/http/v1"
	"project-skbackend/internal/di"
	"project-skbackend/packages/servers"
	"project-skbackend/packages/utils/utlogger"
	"syscall"

	"github.com/gin-gonic/gin"
)

func Run(cfg *configs.Config) {
	// * setup context
	ctx := context.Background()

	// * init config
	i := configs.NewInitConfig(ctx, *cfg)
	i, err := i.InitConfig()
	if err != nil {
		utlogger.Fatal(err)
	}

	// * close the init after finished
	defer i.Close()

	// * setup new dependency injection
	di := di.NewDependencyInjection(i.GormDB, i.Channel, cfg, i.RedisDB, ctx)

	// * setup consumer
	di.InitServices()

	var forever chan struct{}

	// * HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, i.GormDB, cfg, di, i.RedisDB)
	server := servers.NewServer(handler, servers.Port(cfg.HTTP.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		utlogger.Info("app run: " + s.String())
	case err := <-server.Notify():
		utlogger.Fatal(err)
	}

	err = server.Shutdown()
	if err != nil {
		utlogger.Fatal(err)
	}

	<-forever
}
