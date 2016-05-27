package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Token  string `json:"token"`
	Listen string `json:"listen"`

	RedditAgentFile string `json:"reddit_agent_file"`

	Github []GithubConfig `json:"github"`
	Reddit []RedditConfig `json:"reddit"`
}

type GithubConfig struct {
	Repo    string `json:"repo"`
	Channel string `json:"channel"`
	Secret  string `json:"secret"`
}

type RedditConfig struct {
	Sub     string `json:"sub"`
	Channel string `json:"channel"`
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
