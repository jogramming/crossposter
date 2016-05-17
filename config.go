package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Token   string `json:"token"`
	Secret  string `json:"secret"`
	Channel string `json:"channel"`
	Guild   string `json:"guild"`
	Listen  string `json:"listen"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = json.Unmarshal(file, &config)
	return &config, err
}
