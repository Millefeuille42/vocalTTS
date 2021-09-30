package main

import "github.com/bwmarrin/discordgo"

// discordAgent Contains discord's session and message structs
type discordAgent struct {
	session *discordgo.Session
	message *discordgo.MessageCreate
}

type settings struct {
	Id    string
	Admin string
	Vocal string
	Text  string
}
