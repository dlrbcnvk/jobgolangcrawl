package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DB struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Url      string `yaml:"url"`
	} `yaml:"db"`
	Mail struct {
		SmtpHost string `yaml:"smtp_host"`
		SmtpPort string `yaml:"smtp_port"`
		Sender   string `yaml:"sender"`
		Password string `yaml:"password"`
		Receiver string `yaml:"receiver"`
	} `yaml:"mail"`
}

func LoadConfig(env string) (*Config, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	if runeSlice := []rune(dir); runeSlice[len(runeSlice)-1] == '/' {
		dir = string(runeSlice[:len(runeSlice)-1])
	}
	filename := fmt.Sprintf(dir+"/config/config.%s.yml", env)
	log.Printf("Loading config from %s\n", filename)
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
