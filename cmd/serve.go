package cmd

import (
	"fmt"
	"github.com/cjheppell/kondition/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

var (
	serviceConfigFile string
	kubeConfigFile    string

	rootCmd = &cobra.Command{
		Use:   "kondition",
		Short: "Kondition turns k8s deployment availability statuses into HTTP responses for integrating with status checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return serve()
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&serviceConfigFile, "service-config", "", "service configuration file")
	rootCmd.PersistentFlags().StringVar(&kubeConfigFile, "kubeconfig", "", "Kubeconfig file (leave unset if running internally in a cluster)")
	rootCmd.MarkPersistentFlagRequired("service-config")
}

func serve() error {
	prodConfig := zap.NewProductionConfig()
	prodConfig.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	logger, _ := prodConfig.Build()
	defer logger.Sync()
	sugar := logger.Sugar()

	kubeConfigPath, err := getKubeConfigPath()
	if err != nil {
		return err
	}

	err = server.Listen(kubeConfigPath, serviceConfigFile, sugar)
	if err != nil {
		return fmt.Errorf("error starting Kondition server: %s", err)
	}

	return nil
}

func getKubeConfigPath() (string, error) {
	absPath, err := filepath.Abs(kubeConfigFile)
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

	return absPath,nil
}
