package voices

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/tpc3/Kotone-DiVE/lib/config"
	"io"
	"log"
	"net/http"
	"strconv"
)

var (
	AquestalkProxy *aquestalkProxy
)

type aquestalkProxy struct {
	Info VoiceInfo
}

type AquestalkProxyRequest struct {
	VoiceType string `json:"type"`
	Speed     int    `json:"speed"`
	Koe       string `json:"koe"`
}

func init() {
	AquestalkProxy = &aquestalkProxy{
		Info: VoiceInfo{
			Type:             "aquestalk-proxy",
			Format:           "pcm",
			Container:        "wav",
			ReEncodeRequired: true,
			Enabled:          config.CurrentConfig.Voices.AquestalkProxy.Enabled,
		},
	}
	if !config.CurrentConfig.Voices.AquestalkProxy.Enabled {
		log.Print("WARN: aquestalk-proxy is disabled")
		return
	}
}

// Aquestalk-proxy is the proxy program for old version of aquestalk, which only works on win32 systems.
// https://github.com/Na-x4/aquestalk-proxy
// You can just create the http server which runs following to use this backend:
// 1. Convert koe to aquestalk format
// 2. throw the exact same format json to the tcp socket of aquestalk-proxy
// 3. Decode base64 wav from response

func (voiceSource aquestalkProxy) Synth(content string, voice string) ([]byte, error) {
	request, err := json.Marshal(AquestalkProxyRequest{VoiceType: voice, Speed: 100, Koe: content})
	if err != nil {
		return nil, err
	}
	response, err := http.Post(config.CurrentConfig.Voices.AquestalkProxy.Api, "application/json", bytes.NewBuffer(request))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, errors.New("Invalid response from aquestalk-proxy:" + strconv.Itoa(response.StatusCode))
	}
	bin, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return bin, nil
}

func (voiceSource aquestalkProxy) Verify(voice string) error {
	for _, v := range []string{"dvd", "f1", "f2", "imd1", "jgr", "m1", "m2", "r1"} {
		if v == voice {
			return nil
		}
	}
	return errors.New("no such voice")
}

func (voiceSource aquestalkProxy) GetInfo() VoiceInfo {
	return voiceSource.Info
}
