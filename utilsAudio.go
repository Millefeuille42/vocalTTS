package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
)

func readTTSAudio() ([][]byte, error) {
	var buffer = make([][]byte, 0)
	var opusLen int16

	file, err := os.Open("./data/audio.dca")
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return nil, err
	}

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opusLen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return nil, err
			}
			return buffer, nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return nil, err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opusLen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return nil, err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

func queryTTSAudio(agent discordAgent) error {
	ttsParams := url.Values{}
	ttsURL := "http://api.voicerss.org/?"

	ttsParams.Add("key", os.Getenv("TTS_KEY"))
	ttsParams.Add("hl", "fr-fr")
	ttsParams.Add("c", "MP3")
	ttsParams.Add("f", "48khz_16bit_stereo")
	ttsParams.Add("src", agent.message.Member.Nick+" dit: "+isolateEmotes(agent))
	ttsURL += ttsParams.Encode()

	fmt.Println(agent.message.Member.Nick + " dit: " + isolateEmotes(agent))

	response, err := http.Get(ttsURL)
	if err != nil {
		logErrorToChan(agent, err)
		return err
	}
	bytes, err := readHTTPResponse(response)
	if err != nil {
		logErrorToChan(agent, err)
		return err
	}
	err = ioutil.WriteFile("./data/audio.mp3", bytes, 0677)
	if err != nil {
		logErrorToChan(agent, err)
	}
	return err
}

func encodeTTSAudio() error {
	cmd := exec.Command("bash", "-c",
		"ffmpeg -i ./data/audio.mp3 -f s16le -ar 48000 -ac 2 pipe:1 | dca > ./data/audio.dca")
	err := cmd.Run()
	return err
}

func createTTSAudio(agent discordAgent) ([][]byte, error) {
	err := queryTTSAudio(agent)
	if err != nil {
		logErrorToChan(agent, err)
		return nil, err
	}
	err = encodeTTSAudio()
	if err != nil {
		logErrorToChan(agent, err)
		return nil, err
	}
	audioBuffer, err := readTTSAudio()
	if err != nil {
		logErrorToChan(agent, err)
		return nil, err
	}
	return audioBuffer, nil
}
