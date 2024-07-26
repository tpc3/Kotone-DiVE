package voices

import (
	"github.com/tpc3/Kotone-DiVE/lib/config"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	uri = "https://api.voicetext.jp/v1/tts"
)

var (
	VoiceText voiceText
)

type voiceText struct {
	Info     VoiceInfo
	request  *http.Request
	speakers []string
}

func init() {
	VoiceText = voiceText{
		Info: VoiceInfo{
			Type:             "voicetext",
			Format:           "pcm",
			Container:        "wav",
			ReEncodeRequired: false,
			Enabled:          config.CurrentConfig.Voices.VoiceText.Enabled,
		},
	}
	if !config.CurrentConfig.Voices.VoiceText.Enabled {
		log.Print("WARN: VoiceText is disabled")
		return
	}
	VoiceText.speakers = []string{"show", "haruka", "hikari", "takeru", "santa", "bear"}
	request, err := http.NewRequest(http.MethodPost, uri, nil)
	if err != nil {
		log.Fatal(err)
	}
	request.SetBasicAuth(config.CurrentConfig.Voices.VoiceText.Token, "")
	VoiceText.request = request
}

func (voiceSource voiceText) Synth(content string, voice string) ([]byte, error) {
	var text string
	runeContent := []rune(content)
	if len(runeContent) > 200 {
		text = string(runeContent[:200])
	} else {
		text = content
	}

	values := url.Values{}
	values.Set("text", text)
	values.Set("speaker", voice)

	// copy
	req := *voiceSource.request

	r := strings.NewReader(values.Encode())
	req.Body = io.NopCloser(r)
	req.ContentLength = int64(r.Len())
	req.GetBody = func() (io.ReadCloser, error) { return req.Body, nil }
	// If content type is not specified, they will return:
	// 400 {"error":{"message":"speaker must be specified"}}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
	return bin, nil
}

func (voiceSource voiceText) Verify(voice string) error {
	for _, v := range voiceSource.speakers {
		if voice == v {
			return nil
		}
	}
	return errors.New("no such voice")
}

func (voiceSource voiceText) GetInfo() VoiceInfo {
	return voiceSource.Info
}
