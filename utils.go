package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

// logErrorToChan Sends plain error to command channel
func logErrorToChan(agent discordAgent, err error) {
	if err == nil {
		return
	}
	logError(err)
	_, _ = agent.session.ChannelMessageSend(agent.message.ChannelID,
		fmt.Sprintf("An Error occured, Please Try Again Later {%s}", err.Error()))
}

// createDirIfNotExist Check if dir exists, if not create it
func createDirIfNotExist(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(path, os.ModePerm)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

// createFileIfNotExist Check if file exists, if not create it
func createFileIfNotExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		_, err := os.Create(path)
		if err != nil {
			logError(err)
			return false, err
		}
		return false, nil
	}
	return true, nil
}

// logError Prints error + StackTrace to stderr if error
func logError(err error) {
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err, string(debug.Stack()))
	}
}

// checkError Panic if error
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// setUpCloseHandler Set up a handler for Ctrl+C and closing the bot
func setUpCloseHandler(session *discordgo.Session) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		_ = session.Close()
		os.Exit(0)
	}()
}

// readHTTPResponse extracts the body from a http response
func readHTTPResponse(response *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return []byte(""), err
	}
	return body, nil
}
