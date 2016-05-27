package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"io/ioutil"
	"log"
	"net/http"
)

func handleGithub(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RequestURI)
	event := r.Header.Get("X-GitHub-Event")

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	switch event {
	case "push":
		handlePush(data, r)
	}
}

func handlePush(data []byte, r *http.Request) {
	var pushEvent github.PushEvent
	err := json.Unmarshal(data, &pushEvent)
	if err != nil {
		log.Println("Error decoding payload", err)
		return
	}

	repo := pushEvent.Repo
	if repo == nil {
		log.Println("No repo")
		return
	}

	ok := false
	var githubConfig GithubConfig
	for _, v := range config.Github {
		if v.Repo == *repo.Name {
			ok = CheckSignature(data, []byte(v.Secret), r)
			githubConfig = v
		}
	}

	if !ok {
		log.Println("Unknown repo or incorrect secret", *repo.Name)
		return
	}

	pusher := ""
	if pushEvent.Pusher.Name != nil {
		pusher = *pushEvent.Pusher.Name
	}

	header := fmt.Sprintf("**%s** Pushed %d commit(s) to **%s**\n", pusher, len(pushEvent.Commits), *repo.Name)
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

		body += fmt.Sprintf("%s\n**%s**: `%s`\n\n", url, name, msg)
	}

	fullMessage := header + "\n" + body
	log.Println("Sending github message from repo ", *repo.Name)
	_, err = dgo.ChannelMessageSend(githubConfig.Channel, fullMessage)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}

func CheckSignature(data, secret []byte, r *http.Request) bool {
	signature := r.Header.Get("X-Hub-Signature")

	return ValidateHMACDigest(data, signature, secret)
}

// CheckMAC reports whether messageMAC is a valid HMAC tag for message.
func ValidateHMACDigest(message []byte, messageMAC string, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)

	expectedMAC := fmt.Sprintf("sha1=%x", mac.Sum(nil))
	return expectedMAC == messageMAC
}
