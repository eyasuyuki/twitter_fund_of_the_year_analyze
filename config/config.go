package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Year         string `json:"year"`
	PageUrl      string `json:"page_url"`
	TogetterUrl  string `json:"togetter_url"`
	DatabaseName string `json:"database_name"`
	ReportFile   string `json:"report_file"`
}

func NewConfig(configName string) *Config {
	if configName == "" {
		configName = "config.json"
	}
	cf := Config{}
	f, err := os.Open(configName)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	err = json.NewDecoder(f).Decode(&cf)
	if err != nil {
		log.Fatal(err)
	}

	return &cf
}
