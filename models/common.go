package models

func InitSomeThing() {
	// viper is case-insensitive, so all keys in iaas.json should be lowercase
	InitClouds()

	InitDockerClient()
	InitKubernetesClient()
}
