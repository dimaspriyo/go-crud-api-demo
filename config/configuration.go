package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DB struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Name     string `yaml:"name"`
		Config   string `yaml:"config"`
		Port     string `yaml:"port"`
	} `yaml:"database"`
}

func ReadConfig() (config *Config, err error) {

	pwd, _ := os.Getwd()
	configPath := pwd + "/config/database.yml"
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}
	return config, nil

}
