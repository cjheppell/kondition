package main

import (
	"github.com/cjheppell/kondition/server"
	"go.uber.org/zap"
	"os"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	err := server.Listen(sugar)
	if err != nil {
		sugar.Errorf("error starting Kondition server: %s", err)
		os.Exit(1)
	}
}
