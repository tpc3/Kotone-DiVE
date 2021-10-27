package voices

import (
	"Kotone-DiVE/lib/config"
	"context"
	"errors"
	"log"
	"strings"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"google.golang.org/api/option"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

var (
	client *texttospeech.Client
	ctx    context.Context
)

const Gcp = "gcp"

func init() {
	if !config.CurrentConfig.Voices.Gcp.Enabled {
		log.Print("WARN: Gcp is disabled")
		return
	}
	ctx = context.Background()
	var err error
	client, err = texttospeech.NewClient(ctx, option.WithCredentialsJSON([]byte(config.CurrentConfig.Voices.Gcp.Token)))
	if err != nil {
		log.Fatal(err)
	}

}

func GcpClose() {
	if config.CurrentConfig.Voices.Gcp.Enabled {
		client.Close()
	}
}

func GcpSynth(content *string, voice *string) (*[]byte, error) {
	lang, err := gcpLang(voice)
	if err != nil {
		return nil, err
	}
	response, err := client.SynthesizeSpeech(ctx, &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: *content,
			},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: lang,
			Name:         *voice,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_OGG_OPUS,
		},
	})
	if err != nil {
		return nil, err
	}
	c := response.GetAudioContent()
	return &c, nil
}

func gcpLang(voice *string) (string, error) {
	arr := strings.SplitN(*voice, "-", 4)
	if len(arr) != 4 {
		return "", errors.New("can't detect voice lang:" + *voice)
	}
	return arr[0] + "-" + arr[1], nil
}

func GcpVerify(voice *string) error {
	lang, err := gcpLang(voice)
	if err != nil {
		return err
	}
	response, err := client.ListVoices(ctx, &texttospeechpb.ListVoicesRequest{LanguageCode: lang})
	if err != nil {
		return err
	}
	for _, v := range response.GetVoices() {
		if v.GetName() == *voice {
			return nil
		}
	}
	return errors.New("No such voice:" + *voice)
}
