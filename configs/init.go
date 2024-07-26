package configs

import (
	"context"
	"project-skbackend/packages/utils/utlogger"

	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type (
	Init struct {
		ctx context.Context
		cfg Config

		Channel *amqp091.Channel
		GormDB  *gorm.DB
		RedisDB *redis.Client

		clqueue func()
	}
)

func NewInitConfig(ctx context.Context, cfg Config) *Init {
	return &Init{
		ctx: ctx,
		cfg: cfg,
	}
}

func (i *Init) InitConfig() (*Init, error) {
	// * setup database
	db, err := i.initDB()
	if err != nil {
		return nil, err
	}
	i.GormDB = db

	// * setup redis
	rdb := i.initRedis()
	i.RedisDB = rdb

	// * setup queue
	ch, clqueue, err := i.initQueue()
	if err != nil {
		return nil, err
	}
	i.Channel = ch
	i.clqueue = clqueue

	return i, err
}

func (i *Init) initDB() (*gorm.DB, error) {
	// * get the database object
	db, err := gorm.Open(postgres.Open(i.cfg.DB.GetDbConnectionUrl()), &gorm.Config{
		Logger: logger.Default.LogMode(i.cfg.GetLogLevel()),
	})
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// * setup the database migration
	err = i.cfg.DB.DBSetup(db)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return db, nil
}

func (i *Init) initRedis() *redis.Client {
	// * get redis client
	rdb := i.cfg.Redis.GetRedisClient()

	return rdb
}

func (i *Init) initQueue() (*amqp.Channel, func(), error) {
	// * setup rabbit mq
	ch, clqueue := i.cfg.Queue.Init()
	i.cfg.Queue.SetupRabbitMQ(ch, i.cfg)

	return ch, clqueue, nil
}

func (i *Init) Close() {
	if i.clqueue != nil {
		i.clqueue()
	}
}
