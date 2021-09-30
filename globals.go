package main

import (
	"github.com/bwmarrin/discordgo"
	"sync"
)

var speakingLock sync.RWMutex
var commonLock sync.RWMutex
var vocalConnection *discordgo.VoiceConnection
var isConnected bool
