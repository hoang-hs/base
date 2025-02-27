package main

import (
	"context"
	"flag"
	"fmt"
	bootstrap2 "github.com/hoang-hs/base/src/bootstrap"
	"github.com/hoang-hs/base/src/common"
	log2 "github.com/hoang-hs/base/src/common/log"
	"github.com/hoang-hs/base/src/configs"
	"github.com/hoang-hs/base/src/pkg"
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
	log2.NewLogger()
}

func main() {
	logger := log2.GetLogger().GetZap()
	logger.Debugf("App %s is running", configs.Get().Mode)
	app := fx.New(
		fx.Provide(log2.GetLogger().GetZap),
		fx.Invoke(pkg.InitTracer),

		// storage module
		bootstrap2.BuildStorageModules(),

		// build service module
		bootstrap2.BuildServicesModules(),

		// build http server
		bootstrap2.BuildValidator(),
		bootstrap2.BuildControllerModule(),

		// build consumer
		bootstrap2.BuildConsumerModule(),

		bootstrap2.BuildHTTPServerModule(),
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
