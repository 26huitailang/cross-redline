package data

import (
	"database/sql"
	"xredline/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewUserRepo, NewOperationLogRepo, NewCaptchaRepo)

// Data .
type Data struct {
	db    *sql.DB
	redis *redis.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	db, err := sql.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		return nil, nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Network:  c.Redis.Network,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	cleanup := func() {
		if err := db.Close(); err != nil {
			log.NewHelper(logger).Errorf("failed to close database: %v", err)
		}
		if err := redisClient.Close(); err != nil {
			log.NewHelper(logger).Errorf("failed to close redis: %v", err)
		}
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &Data{db: db, redis: redisClient}, cleanup, nil
}
