package voices

import (
	"Kotone-DiVE/lib/config"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	uri       = "https://api.voicetext.jp/v1/tts"
	VoiceText = "voicetext"
)

var (
	speakers  []string
	vtRequest *http.Request
)

func init() {
	if !config.CurrentConfig.Voices.VoiceText.Enabled {
		log.Print("WARN: VoiceText is disabled")
		return
	}
	speakers = []string{"show", "haruka", "hikari", "takeru", "santa", "bear"}
	var err error
	vtRequest, err = http.NewRequest(http.MethodPost, uri, nil)
	if err != nil {
		log.Fatal(err)
	}
	vtRequest.SetBasicAuth(config.CurrentConfig.Voices.VoiceText.Token, "")
}

func VoiceTextSynth(content *string, voice *string) (*[]byte, error) {
	var text string
	runeContent := []rune(*content)
	if len(runeContent) > 200 {
		text = string(runeContent[:200])
	} else {
		text = *content
	}

	values := url.Values{}
	values.Set("text", text)
	values.Set("speaker", *voice)

	// copy
	req := *vtRequest

	r := strings.NewReader(values.Encode())
	req.Body = io.NopCloser(r)
	req.ContentLength = int64(r.Len())
	req.GetBody = func() (io.ReadCloser, error) { return req.Body, nil }
	// If content type is not specified, they will return:
	// 400 {"error":{"message":"speaker must be specified"}}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	log.Printf("%+v", req)
	log.Print(values.Encode())

	res, err := httpCli.Do(&req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bin, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("Response code error from voicetext:" + strconv.Itoa(res.StatusCode) + " " + string(bin))
	}
	return &bin, nil
}

func VoiceTextVerify(voice *string) error {
	for _, v := range speakers {
		if *voice == v {
			return nil
		}
	}
	return errors.New("no such voice")
}
