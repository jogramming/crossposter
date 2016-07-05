package main

import (
	"fmt"
	"github.com/turnage/graw"
	"github.com/turnage/redditproto"
	"log"
)

// autoReplier is a grawbot, which is just a struct that implements methods graw
// looks for, which handle the events from Reddit graw feeds it.
type RedditBot struct {
	eng graw.Engine
}

// SetUp is a method graw looks for. If it is implemented, it will be called
// before the engine starts looking for events on Reddit. If SetUp returns an
// error, the bot will stop.
func (r *RedditBot) SetUp() error {
	r.eng = graw.GetEngine(r)
	log.Println("Reddit Bot is set up!")
	return nil
}

// Called when a post is made
func (r *RedditBot) Post(post *redditproto.Link) {
	channel := ""
	for _, v := range config.Reddit {
		if v.Sub == post.GetSubreddit() {
			channel = v.Channel
			break
		}
	}

	if channel == "" {
		log.Println("Channel for subreddit", post.GetSubreddit(), "not found")
		return
	}

	author := post.GetAuthor()
	sub := post.GetSubreddit()

	typeStr := "link"
	if post.GetIsSelf() {
		typeStr = "self post"
	}

	body := fmt.Sprintf("/u/**%s** Posted a new %s in **/r/%s**:\n**%s**\n%s\n", author, typeStr, sub, post.GetTitle(), "http://reddit.com/"+post.GetPermalink())

	if post.GetIsSelf() {
		body += fmt.Sprintf("*%s*", post.GetSelftext())
	} else {
		body += post.GetUrl() + "\n"
	}

	log.Println("Posting a new reddit message to", sub)
	_, err := dgo.ChannelMessageSend(channel, body)
	if err != nil {
		log.Println("Error posting message", err)
	}
}
