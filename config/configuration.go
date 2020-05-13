package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

//Database Config Struct
type YAMLConfig struct {
	DB struct {
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Name     string `yaml:"name"`
		Config   string `yaml:"config"`
		Port     string `yaml:"port"`
	} `yaml:"database"`
}

//JWT Config Struct
type JSONConfig struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

func ReadYAMLConfig() (config *YAMLConfig, err error) {

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

func ReadJSONConfig() (config JSONConfig, err error) {

	pwd, _ := os.Getwd()
	configPath := pwd + "/config/jwt.json"
	file, err := os.Open(configPath)
	if err != nil {
		return config, err
	}

	read, err := ioutil.ReadAll(file)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(read, &config)
	if err != nil {
		return config, err
	}

	return config, nil

}
