package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Domain string `yaml:"domain"`
	Port   string `yaml:"port"`
}

func GetConfig() *Config {
	config := &Config{}

	configFile, err := os.Open("../config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err2 := yaml.NewDecoder(configFile).Decode(&config); 
	if err2 != nil {
		log.Fatal(err2)
	}

	return config
}