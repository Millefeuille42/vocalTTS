package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func loadFile(agent discordAgent, id string) (settings, error) {
	user := settings{}

	if id == "" {
		id = agent.message.Author.ID
	}

	fileData, err := ioutil.ReadFile(fmt.Sprintf("./data/%s.json", id))
	if err != nil {
		logErrorToChan(agent, err)
		return settings{}, err
	}

	err = json.Unmarshal(fileData, &user)
	if err != nil {
		logErrorToChan(agent, err)
		return settings{}, err
	}

	return user, nil
}

func writeFile(data settings, agent discordAgent) error {
	dataBytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		logErrorToChan(agent, err)
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("./data/%s.json", data.Id), dataBytes, 0677)
	if err != nil {
		logErrorToChan(agent, err)
		return err
	}
	return nil
}
