package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Parse(configPath string) ([]Service, error) {
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	_, err = os.Stat(absConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file '%s' does not exist", configPath)
		}
		return nil, err
	}

	yamlFile, err := ioutil.ReadFile(absConfigPath)
	if err != nil {
		return nil, err
	}

	var serviceConfigs []Service
	err = yaml.Unmarshal(yamlFile, &serviceConfigs)
	if err != nil {
		return nil, err
	}
	applyDefaultsTo(&serviceConfigs)

	return serviceConfigs, nil
}

func applyDefaultsTo(services *[]Service) {
	for i := range *services {
		if (*services)[i].MinReadyReplicas < 1 {
			(*services)[i].MinReadyReplicas = 1
		}
	}
}
