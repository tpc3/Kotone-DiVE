package voices

import (
	"Kotone-DiVE/lib/config"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"strconv"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/watson-developer-cloud/go-sdk/v2/texttospeechv1"
)

var (
	Watson *watson
)

type watson struct {
	Info VoiceInfo
	tts  *texttospeechv1.TextToSpeechV1
}

func init() {
	Watson = &watson{
		Info: VoiceInfo{
			Type:             "watson",
			Format:           "opus",
			Container:        "ogg",
			ReEncodeRequired: false,
			Enabled:          config.CurrentConfig.Voices.Watson.Enabled,
		},
	}
	if !config.CurrentConfig.Voices.Watson.Enabled {
		log.Print("WARN: Watson is disabled")
		return
	}
	auth := &core.IamAuthenticator{ApiKey: config.CurrentConfig.Voices.Watson.Token}
	tts, err := texttospeechv1.NewTextToSpeechV1(&texttospeechv1.TextToSpeechV1Options{Authenticator: auth})
	if err != nil {
		log.Fatal("Watson init error:", err)
	}
	err = tts.SetServiceURL(config.CurrentConfig.Voices.Watson.Api)
	if err != nil {
		log.Print(err)
		return
	}

	Watson.tts = tts
}

func (voiceSource watson) Synth(content string, voice string) ([]byte, error) {
	var buf bytes.Buffer
	err := xml.EscapeText(&buf, []byte(content))
	if err != nil {
		return nil, err
	}
	str := buf.String()
	result, response, err := voiceSource.tts.Synthesize(&texttospeechv1.SynthesizeOptions{
		Text:   &str,
		Accept: core.StringPtr("audio/ogg;codecs=opus"),
		Voice:  &voice,
	})
	if config.CurrentConfig.Debug {
		log.Print(response)
	}
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		// ???
		return nil, errors.New("Invalid status code from Watson:" + strconv.Itoa(response.StatusCode))
	}
	if result != nil {
		bin, err := io.ReadAll(result)
		if err != nil {
			return nil, err
		}
		err = result.Close()
		if err != nil {
			return nil, err
		}
		return bin, nil
	}
	return nil, nil
}

func (voiceSource watson) Verify(voice string) error {
	result, response, err := voiceSource.tts.ListVoices(&texttospeechv1.ListVoicesOptions{})
	if config.CurrentConfig.Debug {
		log.Print(response)
	}
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New("Invalid status code from Watson:" + strconv.Itoa(response.StatusCode))
	}
	if result != nil {
		for _, v := range result.Voices {
			if *v.Name == voice {
				return nil
			}
		}
	}
	return errors.New("Voice is not implemented:" + voice)
}

func (voiceSource watson) GetInfo() VoiceInfo {
	return voiceSource.Info
}
