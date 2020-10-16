package config

type Service struct {
	Name             string `yaml:"name"`
	DeploymentName   string `yaml:"deploymentName"`
	Namespace        string `yaml:"namespace"`
	ApiPath          string `yaml:"apiPath"`
	MinReadyReplicas int32  `yaml:"minReadyReplicas"`
}
