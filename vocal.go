package main

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"time"
)

func joinVocalChannel(agent discordAgent, localSettings settings) (v *discordgo.VoiceConnection, err error) {
	if !isConnected {
		if localSettings.Vocal == "" {
			_, _ = agent.session.ChannelMessageSend(agent.message.ChannelID,
				"You must select a vocal channel first!")
			return nil, errors.New("incomplete setup")
		}
		v, err = agent.session.ChannelVoiceJoin(localSettings.Id, localSettings.Vocal, false, false)
		if err != nil {
			logErrorToChan(agent, err)
			return nil, err
		}
		isConnected = true
		vocalConnection = v
		return
	}
	v = vocalConnection
	return
}

func speak(agent discordAgent, localSettings settings) {
	commonLock.Lock()
	v, err := joinVocalChannel(agent, localSettings)
	if err != nil {
		commonLock.Unlock()
		return
	}
	audioBuffer, err := createTTSAudio(agent)
	if err != nil {
		_ = v.Disconnect()
		isConnected = false
		commonLock.Unlock()
		return
	}
	commonLock.Unlock()

	speakingLock.Lock()
	time.Sleep(250 * time.Millisecond)
	logErrorToChan(agent, v.Speaking(true))
	for _, buff := range audioBuffer {
		v.OpusSend <- buff
	}
	logErrorToChan(agent, v.Speaking(false))
	time.Sleep(250 * time.Millisecond)
	speakingLock.Unlock()
}
