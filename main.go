package main

import (
	"fmt"
	"github.com/cjheppell/kondition/server"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	kubeConfigPath, err := getKubeConfigPath(os.Args)
	if err != nil {
		sugar.Errorf(err.Error())
		os.Exit(1)
	}

	err = server.Listen(kubeConfigPath, sugar)
	if err != nil {
		sugar.Errorf("error starting Kondition server: %s", err)
		os.Exit(1)
	}
}

func getKubeConfigPath(args []string) (string, error) {
	kubeConfigPath := ""
	if len(args) > 1 {
		providedPath := os.Args[1]
		absPath, err := filepath.Abs(providedPath)
		if err != nil {
			return "", fmt.Errorf("error resolving path for provided kubeconfig file. err: %s", err)
		}

		file, err := os.Stat(absPath)
		if err != nil {
			if os.IsNotExist(err) {
				return "", fmt.Errorf("provided kubeconfig file '%s' does not exist", absPath)
			}
		}
		if file.IsDir() {
			return "", fmt.Errorf("provided path '%s' was a directory, not a file", absPath)
		}

		kubeConfigPath = absPath
	}

	return kubeConfigPath, nil
}