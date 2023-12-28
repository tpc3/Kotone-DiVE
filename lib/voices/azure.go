package voices

import (
	"Kotone-DiVE/lib/config"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

const (
	ssml = "<speak version='1.0' xml:lang='%s'><voice xml:lang='%s' xml:gender='%s' name='%s'>%s</voice></speak>"
)

var (
	Azure azure
)

type azure struct {
	Info    VoiceInfo
	voices  map[string]azureVoice
	request *http.Request
}

type azureVoice struct {
	ShortName string
	Gender    string
	Locale    string
}

func init() {
	baseUrl := "https://" + config.CurrentConfig.Voices.Azure.Region + ".tts.speech.microsoft.com/cognitiveservices/"
	Azure = azure{
		Info: VoiceInfo{
			Type:             "azure",
			Format:           "opus",
			Container:        "ogg",
			ReEncodeRequired: false,
			Enabled:          config.CurrentConfig.Voices.Azure.Enabled,
		},
	}
	if !config.CurrentConfig.Voices.Azure.Enabled {
		log.Print("WARN: azure is disabled")
		return
	}
	request, err := http.NewRequest(http.MethodPost, baseUrl+"v1", nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("Ocp-Apim-Subscription-Key", config.CurrentConfig.Voices.Azure.Key)
	request.Header.Add("Content-Type", "application/ssml+xml")
	request.Header.Add("X-Microsoft-OutputFormat", "ogg-48khz-16bit-mono-opus")
	Azure.request = request

	getReq, err := http.NewRequest(http.MethodGet, baseUrl+"voices/list", nil)
	if err != nil {
		log.Fatal(err)
	}
	getReq.Header.Add("Ocp-Apim-Subscription-Key", config.CurrentConfig.Voices.Azure.Key)
	if httpCli == nil {
		httpCli = &http.Client{}
	}
	res, err := httpCli.Do(getReq)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if config.CurrentConfig.Debug {
		log.Print(res)
	}
	var (
		voices []azureVoice
		bin    []byte
	)
	bin, err = io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bin, &voices)
	if err != nil {
		log.Fatal(err)
	}
	Azure.voices = map[string]azureVoice{}
	for _, v := range voices {
		Azure.voices[v.ShortName] = v
	}
}

func (voiceSource azure) Synth(content string, voice string) ([]byte, error) {
	val, exists := voiceSource.voices[voice]
	if !exists {
		return nil, errors.New("invalid voice type")
	}
	buffer := bytes.Buffer{}
	err := xml.EscapeText(&buffer, []byte(content))
	if err != nil {
		return nil, err
	}
	tmpReq := voiceSource.request
	tmpReq.Body = io.NopCloser(strings.NewReader(fmt.Sprintf(ssml, val.Locale, val.Locale, val.Gender, val.ShortName, buffer.String())))
	cli := http.Client{}
	res, err := cli.Do(tmpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bin, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bin, nil
}

func (voiceSource azure) Verify(voice string) error {
	_, exists := voiceSource.voices[voice]
	if !exists {
		return errors.New("invalid voice type: " + voice)
	}
	return nil
}

func (voiceSource azure) GetInfo() VoiceInfo {
	return voiceSource.Info
}
