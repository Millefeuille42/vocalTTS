package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)

func checkSpeakValidity(agent discordAgent) {
	localSettings, err := loadFile(agent, agent.message.GuildID)
	if err != nil {
		_, _ = agent.session.ChannelMessageSend(agent.message.ChannelID,
			"As the error said, you must init this server!")
		return
	}
	if agent.message.ChannelID == localSettings.Text {
		speak(agent, localSettings)
		return
	}
}

func sendHelp(agent discordAgent) {
	_, _ = agent.session.ChannelMessageSend(agent.message.ChannelID,
		"`!set-text` To set the channel you sent this command on to be the text channel\n"+
			"`!set-vocal <VOCAL ID>` To set the vocal channel\n"+
			"`!leave` To make the bot leave the channel he is currently in\n"+
			"`!help` To print this message!")
}

func leaveChat() {
	if !isConnected {
		return
	}
	commonLock.Lock()
	speakingLock.Lock()
	_ = vocalConnection.Speaking(false)
	_ = vocalConnection.Disconnect()
	speakingLock.Unlock()
	isConnected = false
	commonLock.Unlock()
}

// commandRouter Router for admin
func commandRouter(agent discordAgent) {
	switch {
	case agent.message.Content == "!help":
		sendHelp(agent)
		return
	case strings.HasPrefix(agent.message.Content, "!set-text"):
		settingsSet(agent, false)
		return
	case strings.HasPrefix(agent.message.Content, "!set-vocal"):
		settingsSet(agent, true)
		return
	case agent.message.Content == "!leave":
		leaveChat()
		return
	}
	checkSpeakValidity(agent)
}

// messageHandler Discord bot message handler
func messageHandler(session *discordgo.Session, message *discordgo.MessageCreate) {
	botID, _ := session.User("@me")
	agent := discordAgent{
		session: session,
		message: message,
	}

	if message.Author.ID == botID.ID {
		return
	}
	commandRouter(agent)
}
