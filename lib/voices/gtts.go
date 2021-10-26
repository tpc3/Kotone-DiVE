package voices

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/evalphobia/google-tts-go/googletts"
)

const (
	Gtts  = "gtts"
	token = "368668.249914"
)

func GttsSynth(content *string, voice *string) (*[]byte, error) {
	url, err := googletts.GetTTSURLWithOption(googletts.Option{
		Lang:  *voice,
		Token: token,
		Text:  *content,
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.New("Invalid response from gtts:" + strconv.Itoa(response.StatusCode))
	}
	bin, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &bin, nil
}

func GttsVerify(voice *string) error {
	str := "test"
	_, err := GttsSynth(&str, voice)
	return err
}
