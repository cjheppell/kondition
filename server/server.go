package server

import (
	"fmt"
	"github.com/cjheppell/kondition/config"
	"github.com/cjheppell/kondition/internal/kubernetes"
	"go.uber.org/zap"
	"net/http"
)

func defaultHandler(w http.ResponseWriter, r *http.Request, logger *zap.SugaredLogger) {
	w.WriteHeader(404)
	_, err := fmt.Fprint(w, "Requested service not found")
	if err != nil {
		logger.Warnf("error writing to http response: %s", err)
	}
}

func Listen(kubeConfigPath, serviceConfigPath string, logger *zap.SugaredLogger) error {
	statusClient, err := kubernetes.NewStatusClient(kubeConfigPath, logger)
	if err != nil {
		return err
	}

	services, err := config.Parse(serviceConfigPath)
	if err != nil {
		return err
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { defaultHandler(w, r, logger) })
	err = registerTargets(services, statusClient, logger)
	if err != nil {
		return err
	}
	logger.Info("Kondition is live and listening on port 8080")
	return http.ListenAndServe(":8080", nil)
}

func registerTargets(services []config.Service, statusClient *kubernetes.StatusClient, logger *zap.SugaredLogger) error {
	for _, service := range services {
		logger.Infof("Starting service registration for %s at path %s", service.Name, service.ApiPath)
		registerTarget(service, statusClient, logger)
	}

	return nil
}

func registerTarget(service config.Service, statusClient *kubernetes.StatusClient, logger *zap.SugaredLogger) {
	http.HandleFunc(fmt.Sprintf("%s", service.ApiPath), func(w http.ResponseWriter, r *http.Request) {
		logger.Debugf("request URL was: %s", r.URL)
		isReady, err := statusClient.IsDeploymentReady(service.DeploymentName, service.Namespace)
		if err != nil {
			logger.Errorf("error getting service status for service '%s'. err: %s", service.Name, err)
		}

		if !isReady {
			w.WriteHeader(503)
			_, err := fmt.Fprintf(w, "Service %s is unavailable", service.Name)
			if err != nil {
				logger.Errorf("error writing unavailable status to http response. err: %s", err)
			}
		} else {
			w.WriteHeader(200)
			_, err := fmt.Fprintf(w, "Service %s is available", service.Name)
			if err != nil {
				logger.Errorf("error writing available status to http response. err: %s", err)
			}
		}
	})
}