package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type S3Config struct {
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	Region    string `yaml:"region"`
	UseSSL    bool   `yaml:"useSSL"`
}

type Config struct {
	S3 S3Config `yaml:"s3"`
}

var Cfg Config

func LoadConfig() {
	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	err = yaml.Unmarshal(file, &Cfg)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err)
	}
}
