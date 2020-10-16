package server

import (
	"fmt"
	"github.com/cjheppell/kondition/config"
	"github.com/cjheppell/kondition/internal/kubernetes"
	"go.uber.org/zap"
	"net/http"
)

func defaultHandler(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger) {
	_, err := fmt.Fprint(w, "Kondition is live and running. Navigate to a path to check the status of a watched service.")
	if err != nil {
		logger.Warnf("error writing to http response: %s", err)
	}
}

func Listen(kubeConfigPath, serviceConfigPath string, logger *zap.SugaredLogger) error {
	_, err := kubernetes.NewStatusClient(kubeConfigPath)
	if err != nil {
		return err
	}

	_, err = config.Parse(serviceConfigPath)
	if err != nil {
		return err
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { defaultHandler(w, r, logger) })
	logger.Info("Kondition is live and listening on port 8080")
	return http.ListenAndServe(":8080", nil)
}
