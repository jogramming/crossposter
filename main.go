package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/turnage/graw"
	"log"
	"net/http"
)

var (
	config *Config
	dgo    *discordgo.Session
)

func main() {

	var err error
	config, err = LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	dgo, err = discordgo.New(config.Token)
	if err != nil {
		panic(err)
	}
	err = dgo.Open()
	if err != nil {
		panic(err)
	}
	go RunReddit()
	go RunGithub()
	select {}
}

func RunGithub() {
	if len(config.Github) < 1 {
		log.Println("No github sources defined, not running webhook server..")
		return
	}
	log.Println("Starting github webhook server")
	http.HandleFunc("/", handleGithub)
	log.Println(http.ListenAndServe(config.Listen, nil))
}

func RunReddit() {
	if len(config.Reddit) < 1 {
		log.Println("No reddit sources defined, not running reddit bot...")
		return
	}
	agentFile := config.RedditAgentFile
	if config.RedditAgentFile == "" {
		log.Println("No agent file specified, using reddit.agent")
		agentFile = "reddit.agent"
	}
	bot := &RedditBot{}

	subs := make([]string, 0)
	for _, v := range config.Reddit {
		subs = append(subs, v.Sub)
	}
	log.Println("Running graw on ", subs)
	err := graw.Run(agentFile, bot, subs...)
	if err != nil {
		log.Println("Error running graw:", err)
	}
}
