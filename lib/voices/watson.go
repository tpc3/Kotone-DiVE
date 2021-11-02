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
	tts *texttospeechv1.TextToSpeechV1
)

const Watson = "watson"

func init() {
	if !config.CurrentConfig.Voices.Watson.Enabled {
		log.Print("WARN: Voice \"Watson\" is disabled")
		return
	}
	auth := &core.IamAuthenticator{ApiKey: config.CurrentConfig.Voices.Watson.Token}
	var err error
	tts, err = texttospeechv1.NewTextToSpeechV1(&texttospeechv1.TextToSpeechV1Options{Authenticator: auth})
	if err != nil {
		log.Fatal("Watson init error:", err)
	}
	tts.SetServiceURL(config.CurrentConfig.Voices.Watson.Api)
}

func WatsonSynth(content *string, voice *string) (*[]byte, error) {
	var buf bytes.Buffer
	err := xml.EscapeText(&buf, []byte(*content))
	if err != nil {
		return nil, err
	}
	str := buf.String()
	result, response, err := tts.Synthesize(&texttospeechv1.SynthesizeOptions{
		Text:   &str,
		Accept: core.StringPtr("audio/ogg;codecs=opus"),
		Voice:  voice,
	})
	if config.CurrentConfig.Debug {
		log.Print(response)
	}
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		// ???
		return nil, errors.New("Invalid statuscode from Watson:" + strconv.Itoa(response.StatusCode))
	}
	if result != nil {
		bin, err := io.ReadAll(result)
		if err != nil {
			return nil, err
		}
		result.Close()
		return &bin, nil
	}
	return nil, nil
}

func WatsonVerify(voice *string) error {
	result, response, err := tts.ListVoices(&texttospeechv1.ListVoicesOptions{})
	if config.CurrentConfig.Debug {
		log.Print(response)
	}
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return errors.New("Invalid statuscode from Watson:" + strconv.Itoa(response.StatusCode))
	}
	if result != nil {
		for _, v := range result.Voices {
			if *v.Name == *voice {
				return nil
			}
		}
	}
	return errors.New("Voice is not implemented:" + *voice)
}
