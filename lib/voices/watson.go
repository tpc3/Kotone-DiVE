package voices

import (
	"Kotone-DiVE/lib/config"
	"errors"
	"io"
	"log"
	"strconv"

	"github.com/IBM/go-sdk-core/core"
	"github.com/watson-developer-cloud/go-sdk/texttospeechv1"
)

var (
	tts *texttospeechv1.TextToSpeechV1
)

func init() {
	auth := &core.IamAuthenticator{ApiKey: config.CurrentConfig.Voices.Watson.Token}
	var err error
	tts, err = texttospeechv1.NewTextToSpeechV1(&texttospeechv1.TextToSpeechV1Options{Authenticator: auth})
	if err != nil {
		log.Fatal("Watson init error:", err)
	}
	tts.SetServiceURL(config.CurrentConfig.Voices.Watson.Api)
}

func Watson(content *string, voice *string) (*[]byte, error) {
	result, response, err := tts.Synthesize(&texttospeechv1.SynthesizeOptions{
		Text:   content,
		Accept: core.StringPtr("audio/ogg;codecs=opus"),
		Voice:  voice,
	})
	if config.CurrentConfig.Debug {
		log.Print(response)
	}
	if err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		if response.StatusCode != 200 {
			// ???
			return nil, errors.New("Invalid statuscode from Watson:" + strconv.Itoa(response.StatusCode))
		}
		if result != nil {
			bin, err := io.ReadAll(result)
			if err != nil {
				log.Fatal("Watson readall failed:", err)
			}
			result.Close()
			return &bin, nil
		}
		return nil, nil
	}
}