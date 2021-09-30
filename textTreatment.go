package main

import (
	"regexp"
	"strings"
)

func getUserMention(agent discordAgent) (ret string) {
	ret = agent.message.Content
	return
}

func isolateEmotes(agent discordAgent) (ret string) {
	ret = agent.message.Content

	r, _ := regexp.Compile("<:([^:]*):\\d*>")
	emotes := r.FindAllStringSubmatch(agent.message.Content, -1)
	for _, emote := range emotes {
		ret = strings.Replace(ret, emote[0], emote[1], 1)
	}
	return
}
