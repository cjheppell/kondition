package config

import (
	"testing"
)

func TestParsingDefaultMinReadyReplicas(t *testing.T) {
	configFile := "testdata/service1.yaml"
	services, err := Parse(configFile)

	if err != nil {
		t.Errorf("Failed to parse config file: %s", configFile)
	}

	minReadyReplicas := services[0].MinReadyReplicas

	if minReadyReplicas != 1 {
		t.Errorf("Expected minReadyReplicas to be %d, but found %d", 1, minReadyReplicas)
	}
}
