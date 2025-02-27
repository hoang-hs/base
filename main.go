package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/hoang-hs/base/src/bootstrap"
	"github.com/hoang-hs/base/src/common"
	"github.com/hoang-hs/base/src/common/log"
	"github.com/hoang-hs/base/src/configs"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultGracefulTimeout = 15 * time.Second
)

func init() {
	var pathConfig string
	flag.StringVar(&pathConfig, "config", "configs/config.yaml", "path to config file")
	flag.Parse()
	err := configs.LoadConfig(pathConfig)
	if err != nil {
		panic(err)
	}
	if !common.IsProdEnv {
		fmt.Println(configs.Get())
	}
	log.NewLogger()
}

func main() {
	logger := log.GetLogger().GetZap()
	logger.Debugf("App %s is running", configs.Get().Mode)
	app := fx.New(
		fx.Provide(log.GetLogger().GetZap),

		fx.Provide(bootstrap.BuildObserveServicesModules()),

		// storage module
		bootstrap.BuildStorageModules(),

		// build service module
		bootstrap.BuildServicesModules(),

		// build http server
		bootstrap.BuildValidator(),
		bootstrap.BuildControllerModule(),

		// build consumer
		bootstrap.BuildConsumerModule(),

		bootstrap.BuildHTTPServerModule(),
		// bootstrap.BuildGrpcModules(),
	)

	startCtx, cancel := context.WithTimeout(context.Background(), defaultGracefulTimeout)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		logger.Fatalf(err.Error())
	}

	interruptHandle(app, logger)
}

func interruptHandle(app *fx.App, logger *zap.SugaredLogger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Debugf("Listening Signal...")
	s := <-c
	logger.Infof("Received signal: %s. Shutting down Server ...", s)

	stopCtx, cancel := context.WithTimeout(context.Background(), defaultGracefulTimeout)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		logger.Fatalf(err.Error())
	}
}
