package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Port               int    `json:"port"`
	CertFile           string `json:"certFile"`
	CertKey            string `json:"certKey"`
	IsFileServer       bool   `json:"isFileServer"`
	FileServerRootPath string `json:"fileServerRootPath"`
	SecretUserName     string `json:"secretUserName"`
	SecretPassword     string `json:"secretPassword"`
	HmacSampleSecret   string `json:"hmacSampleSecret"`
}

func ReadConfig() (Config, error) {
	data, err := os.ReadFile("conf.json")
	if err != nil {
		return Config{}, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return Config{}, err
	}

	return config, nil
}
