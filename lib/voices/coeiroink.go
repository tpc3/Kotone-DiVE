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
	Coeiroink = "coeiroink"
)

var (
	ciSpeakers         Speakers
	ciSynthesisRequest *http.Request
)

// These structs are currently used from voicevox implementation
// type Speakers struct
// type Speaker struct
// type Style struct

// Coeiroink is the another engine implementation based on the voicevox.
// They uses almost identical api, but requires windows or mac to run. (OSs like Linux requires additional compat layers like wine)

func init() {
	if !config.CurrentConfig.Voices.Coeiroink.Enabled {
		log.Print("WARN: Coeiroink is disabled")
		return
	}
	res, err := http.Get(config.CurrentConfig.Voices.Coeiroink.Api + "/speakers")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	bin, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal([]byte("{ \"speakers\": "+string(bin)+" }"), &ciSpeakers)
	if err != nil {
		log.Fatal(err)
	}
	ciSynthesisRequest, err = http.NewRequest(http.MethodPost, config.CurrentConfig.Voices.Coeiroink.Api+"/synthesis", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func CoeiroinkSynth(content *string, voice *string) (*[]byte, error) {
	id := -1
	for _, speaker := range ciSpeakers.Speakers {
		for _, v := range speaker.Styles {
			if speaker.Name+v.Name == *voice {
				id = v.Id
				break
			}
		}
	}
	if id == -1 {
		return nil, errors.New("no such voice")
	}

	// copy
	res, err := http.Post(config.CurrentConfig.Voices.Coeiroink.Api+"/audio_query?speaker="+strconv.Itoa(id)+"&text="+url.QueryEscape(*content), "", nil)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("Response code error from coeiroink:" + strconv.Itoa(res.StatusCode))
	}

	// copy
	req := *ciSynthesisRequest

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
		return nil, errors.New("Response code error from coeiroink:" + strconv.Itoa(res.StatusCode) + " " + string(bin))
	}
	return &bin, nil
}

func CoeiroinkVerify(voice *string) error {
	for _, speaker := range ciSpeakers.Speakers {
		for _, v := range speaker.Styles {
			if speaker.Name+v.Name == *voice {
				return nil
			}
		}
	}
	return errors.New("no such voice")
}
