package voices

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/tpc3/Kotone-DiVE/lib/config"
)

var AivisSpeech aivisspeech

type aivisspeech struct {
	Info     VoiceInfo
	Speakers Speakers
	Request  *http.Request
}

// These structs are currently used from voicevox implementation
// type Speakers struct
// type Speaker struct
// type Style struct

// AivisSpeech is the another engine implementation based on the voicevox.
// They uses almost identical api, but does not provide official docker image (There's Dockerfile though).

func init() {
	AivisSpeech = aivisspeech{
		Info: VoiceInfo{
			Type:             "aivisspeech",
			Format:           "pcm",
			Container:        "wav",
			ReEncodeRequired: true,
			Enabled:          config.CurrentConfig.Voices.AivisSpeech.Enabled,
		},
	}
	if !config.CurrentConfig.Voices.AivisSpeech.Enabled {
		log.Print("WARN: AivisSpeech is disabled")
		return
	}
	res, err := http.Get(config.CurrentConfig.Voices.AivisSpeech.Api + "/speakers")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	bin, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal([]byte("{ \"speakers\": "+string(bin)+" }"), &AivisSpeech.Speakers)
	if err != nil {
		log.Fatal(err)
	}
	request, err := http.NewRequest(http.MethodPost, config.CurrentConfig.Voices.AivisSpeech.Api+"/synthesis", nil)
	if err != nil {
		log.Fatal(err)
	}
	AivisSpeech.Request = request
}

func (voiceSource aivisspeech) Synth(content string, voice string) ([]byte, error) {
	id := -1
	for _, speaker := range voiceSource.Speakers.Speakers {
		for _, v := range speaker.Styles {
			if speaker.Name+v.Name == voice {
				id = v.Id
				break
			}
		}
	}
	if id == -1 {
		return nil, errors.New("no such voice")
	}

	// copy
	res, err := http.Post(config.CurrentConfig.Voices.AivisSpeech.Api+"/audio_query?speaker="+strconv.Itoa(id)+"&text="+url.QueryEscape(content), "", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("Response code error from aivisspeech:" + strconv.Itoa(res.StatusCode))
	}

	// copy
	req := *voiceSource.Request

	query := res.Body
	buf := new(bytes.Buffer)
	length, err := buf.ReadFrom(query)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = "speaker=" + strconv.Itoa(id)
	req.Body = io.NopCloser(buf)
	req.ContentLength = length
	req.GetBody = func() (io.ReadCloser, error) { return req.Body, nil }
	req.Header.Set("Content-Type", "application/json")

	res, err = httpCli.Do(&req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bin, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("Response code error from aivisspeech:" + strconv.Itoa(res.StatusCode) + " " + string(bin))
	}
	return bin, nil
}

func (voiceSource aivisspeech) Verify(voice string) error {
	for _, speaker := range voiceSource.Speakers.Speakers {
		for _, v := range speaker.Styles {
			if speaker.Name+v.Name == voice {
				return nil
			}
		}
	}
	return errors.New("no such voice")
}

func (voiceSource aivisspeech) GetInfo() VoiceInfo {
	return voiceSource.Info
}
