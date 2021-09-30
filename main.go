package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"time"
)

// startBot Starts discord bot
func startBot() {
	discordBot, err := discordgo.New("Bot " + os.Getenv("BOT_TOKEN"))
	checkError(err)
	discordBot.AddHandler(messageHandler)
	discordBot.AddHandler(readyHandler)
	err = discordBot.Open()
	checkError(err)
	fmt.Println("Discord bot created")

	setUpCloseHandler(discordBot)
}

func readyHandler(s *discordgo.Session, event *discordgo.Ready) {
	_ = s.UpdateListeningStatus("!help")
}

// prepFileSystem Create required directories
func prepFileSystem() error {
	err := createDirIfNotExist("./data")
	return err
}

func main() {
	checkError(prepFileSystem())
	startBot()

	for {
		time.Sleep(time.Second * 3)
	}
}
