package main

import (
	"os"
	"strings"
)

func settingsSet(agent discordAgent, isVocal bool) {
	localSettings := settings{}

	exists, err := createFileIfNotExist("./data/" + agent.message.GuildID + ".json")
	if err != nil {
		return
	}
	if exists {
		localSettings, err = loadFile(agent, agent.message.GuildID)
		if err != nil {
			_ = os.Remove("./data/" + agent.message.GuildID + ".json")
			settingsSet(agent, isVocal)
			return
		}
	} else {
		localSettings.Admin = agent.message.Author.ID
	}
	localSettings.Id = agent.message.GuildID
	if agent.message.Author.ID != localSettings.Admin {
		_, _ = agent.session.ChannelMessageSend(agent.message.ChannelID, "You are not the admin")
		return
	}
	if isVocal {
		settingsSetVocal(agent, localSettings)
		return
	}
	settingsSetText(agent, localSettings)
}

func settingsSetText(agent discordAgent, localSettings settings) {
	localSettings.Text = agent.message.ChannelID
	err := writeFile(localSettings, agent)
	if err != nil {
		return
	}
	_, _ = agent.session.ChannelMessageSend(agent.message.ChannelID, "Success!")
}

func settingsSetVocal(agent discordAgent, localSettings settings) {
	args := strings.Split(agent.message.Content, " ")
	if len(args) <= 1 {
		_, _ = agent.session.ChannelMessageSend(agent.message.ChannelID, "You must provide a channel ID")
	}
	localSettings.Vocal = args[1]
	err := writeFile(localSettings, agent)
	if err != nil {
		return
	}
	_, _ = agent.session.ChannelMessageSend(agent.message.ChannelID, "Success!")
}
