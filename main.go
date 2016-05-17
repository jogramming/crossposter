package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/jonas747/discordgo"
	"io/ioutil"
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

	http.HandleFunc("/", handleGithub)
	log.Println(http.ListenAndServe(config.Listen, nil))
}

func handleGithub(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	event := r.Header.Get("X-GitHub-Event")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	signature := r.Header.Get("X-Hub-Signature")

	if !ValidateHMACDigest(data, signature, []byte(config.Secret)) {
		log.Println("Invalid signature...")
		return
	} else {
		log.Println("Valid signature!")
	}

	switch event {
	case "push":
		handlePush(data)
	}
}

func handlePush(data []byte) {
	var pushEvent github.PushEvent
	err := json.Unmarshal(data, &pushEvent)

	if err != nil {
		panic(err)
	}

	pusher := ""
	if pushEvent.Pusher.Name != nil {
		pusher = *pushEvent.Pusher.Name
	}

	header := fmt.Sprintf("**%s** Pushed %d commit(s)", pusher, len(pushEvent.Commits))
	body := ""
	for _, v := range pushEvent.Commits {
		url := ""
		if v.URL != nil {
			url = *v.URL
		}

		name := ""
		if v.Author.Login != nil {
			name = *v.Author.Login
		} else if v.Author.Name != nil {
			name = *v.Author.Name
		}

		msg := ""
		if v.Message != nil {
			msg = *v.Message
		}

		body += fmt.Sprintf("%s\n**%s**: %s\n", url, name, msg)
	}

	fullMessage := header + "\n" + body
	log.Println("Sending message", fullMessage)
	_, err = dgo.ChannelMessageSend(config.Channel, fullMessage)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}

// CheckMAC reports whether messageMAC is a valid HMAC tag for message.
func ValidateHMACDigest(message []byte, messageMAC string, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)

	expectedMAC := fmt.Sprintf("sha1=%x", mac.Sum(nil))
	return expectedMAC == messageMAC
}
