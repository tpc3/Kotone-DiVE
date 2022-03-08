package voices

import (
	"Kotone-DiVE/lib/config"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	Voicevox = "voicevox"
)

var (
	vvSpeakers         Speakers
	vvSynthesisRequest *http.Request
)

type Speakers struct {
	Speakers []Speaker
}

type Speaker struct {
	Name        string
	SpeakerUUID string
	Styles      []Style
}

type Style struct {
	Id   int
	Name string
}

func init() {
	if !config.CurrentConfig.Voices.Voicevox.Enabled {
		log.Print("WARN: Voicevox is disabled")
		return
	}
	res, err := http.Get(config.CurrentConfig.Voices.Voicevox.Api + "/speakers")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	bin, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal([]byte("{ \"speakers\": "+string(bin)+" }"), &vvSpeakers)
	if err != nil {
		log.Fatal(err)
	}
	vvSynthesisRequest, err = http.NewRequest(http.MethodPost, config.CurrentConfig.Voices.Voicevox.Api+"/synthesis", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func VoicevoxSynth(content *string, voice *string) (*[]byte, error) {
	id := -1
	for _, speaker := range vvSpeakers.Speakers {
		for i, v := range speaker.Styles {
			if speaker.Name+v.Name == *voice {
				id = i
			}
		}
	}
	if id == -1 {
		return nil, errors.New("no such voice")
	}

	// copy
	res, err := http.Post(config.CurrentConfig.Voices.Voicevox.Api+"/audio_query?speaker="+strconv.Itoa(id)+"&text="+url.QueryEscape(*content), "", nil)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("Response code error from voicevox:" + strconv.Itoa(res.StatusCode))
	}

	// copy
	req := *vvSynthesisRequest

	query := res.Body
	buf := new(bytes.Buffer)
	len, err := buf.ReadFrom(query)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = "speaker=" + strconv.Itoa(id)
	req.Body = io.NopCloser(buf)
	req.ContentLength = len
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
		return nil, errors.New("Response code error from voicevox:" + strconv.Itoa(res.StatusCode) + " " + string(bin))
	}
	return &bin, nil
}

func VoicevoxVerify(voice *string) error {
	for _, speaker := range vvSpeakers.Speakers {
		for _, v := range speaker.Styles {
			if speaker.Name+v.Name == *voice {
				return nil
			}
		}
	}
	return errors.New("no such voice")
}
