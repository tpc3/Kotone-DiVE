package voices

import (
	"github.com/tpc3/Kotone-DiVE/lib/config"
	"context"
	"errors"
	"log"
	"strings"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"google.golang.org/api/option"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

var (
	Gcp gcp
)

type gcp struct {
	Info   VoiceInfo
	client *texttospeech.Client
	ctx    context.Context
}

func init() {
	if !config.CurrentConfig.Voices.Gcp.Enabled {
		log.Print("WARN: Gcp is disabled")
		return
	}
	Gcp = gcp{
		Info: VoiceInfo{
			Type:             "gcp",
			Format:           "opus",
			Container:        "ogg",
			ReEncodeRequired: false,
			Enabled:          config.CurrentConfig.Voices.Gcp.Enabled,
		},
	}
	Gcp.ctx = context.Background()
	client, err := texttospeech.NewClient(Gcp.ctx, option.WithCredentialsJSON([]byte(config.CurrentConfig.Voices.Gcp.Token)))
	if err != nil {
		log.Fatal(err)
	}
	Gcp.client = client

}

func (voiceSource gcp) Close() {
	if voiceSource.Info.Enabled {
		voiceSource.client.Close()
	}
}

func (voiceSource gcp) Synth(content string, voice string) ([]byte, error) {
	lang, err := voiceSource.lang(voice)
	if err != nil {
		return nil, err
	}
	response, err := voiceSource.client.SynthesizeSpeech(voiceSource.ctx, &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: content,
			},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: lang,
			Name:         voice,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_OGG_OPUS,
		},
	})
	if err != nil {
		return nil, err
	}
	c := response.GetAudioContent()
	return c, nil
}

func (voiceSource gcp) lang(voice string) (string, error) {
	arr := strings.SplitN(voice, "-", 4)
	if len(arr) != 4 {
		return "", errors.New("can't detect voice lang:" + voice)
	}
	return arr[0] + "-" + arr[1], nil
}

func (voiceSource gcp) Verify(voice string) error {
	lang, err := voiceSource.lang(voice)
	if err != nil {
		return err
	}
	response, err := voiceSource.client.ListVoices(voiceSource.ctx, &texttospeechpb.ListVoicesRequest{LanguageCode: lang})
	if err != nil {
		return err
	}
	for _, v := range response.GetVoices() {
		if v.GetName() == voice {
			return nil
		}
	}
	return errors.New("No such voice:" + voice)
}

func (voiceSource gcp) GetInfo() VoiceInfo {
	return voiceSource.Info
}
