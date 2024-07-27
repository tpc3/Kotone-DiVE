package voices

import (
	"errors"
	"github.com/tpc3/Kotone-DiVE/lib/config"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/evalphobia/google-tts-go/googletts"
)

const (
	token = "368668.249914"
)

var Gtts *gtts

type gtts struct {
	Info VoiceInfo
}

func init() {
	Gtts = &gtts{
		Info: VoiceInfo{
			Type:             "gtts",
			Format:           "mp3",
			Container:        "mp3",
			ReEncodeRequired: false,
			Enabled:          config.CurrentConfig.Voices.Gtts.Enabled,
		},
	}
}

func (voiceSource gtts) Synth(content string, voice string) ([]byte, error) {
	url, err := googletts.GetTTSURLWithOption(googletts.Option{
		Lang:  voice,
		Token: token,
		Text:  content,
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
	if config.CurrentConfig.Debug {
		log.Print(response)
	}
	bin, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return bin, nil
}

func (voiceSource gtts) Verify(voice string) error {
	str := "test"
	_, err := Gtts.Synth(str, voice)
	return err
}

func (voiceSource gtts) GetInfo() VoiceInfo {
	return voiceSource.Info
}
