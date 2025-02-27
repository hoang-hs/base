package bootstrap

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/hoang-hs/base/common"
	"github.com/hoang-hs/base/common/log"
	"github.com/hoang-hs/base/configs"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func BuildStorageModules() fx.Option {
	return fx.Options(
		fx.Provide(newPostgresqlDB),
		fx.Provide(newCacheRedis),
		fx.Provide(newMongoDB),
	)
}

func newPostgresqlDB(lc fx.Lifecycle, config *configs.Config) *gorm.DB {
	cfPostgresql := config.Postgresql
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", cfPostgresql.Host,
		cfPostgresql.Port, cfPostgresql.User, cfPostgresql.DbName, cfPostgresql.SslMode, cfPostgresql.Password)
	logMode := logger.Info
	if common.IsProdEnv {
		logMode = logger.Silent
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		panic(err)
	}
	if config.Tracer.Enabled {
		if err := db.Use(otelgorm.NewPlugin()); err != nil {
			panic(err)
		}
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Debug("Coming OnStop Storage")
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			return sqlDB.Close()
		},
	})
	return db
}

func newCacheRedis(config *configs.Config) redis.UniversalClient {
	cf := config.Redis
	hosts := cf.Hosts
	var client redis.UniversalClient
	isClusterMode := len(hosts) > 1
	if isClusterMode {
		client = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    hosts,
			Username: cf.Username,
			Password: cf.Password,
		})
	} else {
		client = redis.NewClient(&redis.Options{
			Addr:     hosts[0],
			Username: cf.Username,
			Password: cf.Password,
		})
	}

	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.GetLogger().GetZap().Fatalf("ping redis error, err:[%s]", err.Error())
	}
	return client
}

func newMongoDB(lc fx.Lifecycle, cf *configs.Config) *mongo.Database {
	log.Debug("Coming Create Storage")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.ClientOptions{}
	if cf.Tracer.Enabled {
		opts.Monitor = otelmongo.NewMonitor()
	}
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cf.Mongo.Uri), &opts)
	if err != nil {
		log.Fatal("connect mongo db error:[%s]", err.Error())
	}
	if err = client.Ping(context.Background(), nil); err != nil {
		log.Fatal("ping mongo db error:[%s]", err.Error())
	}
	db := client.Database(cf.Mongo.DB)
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info("Coming OnStop Storage")
			return client.Disconnect(ctx)
		},
	})
	return db
}
