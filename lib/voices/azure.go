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
	Azure = "azure"
	ssml  = "<speak version='1.0' xml:lang='%s'><voice xml:lang='%s' xml:gender='%s' name='%s'>%s</voice></speak>"
)

var (
	request     *http.Request
	azureVoices map[string]azureVoice
)

type azureVoice struct {
	ShortName string
	Gender    string
	Locale    string
}

func init() {
	baseUrl := "https://" + config.CurrentConfig.Voices.Azure.Region + ".tts.speech.microsoft.com/cognitiveservices/"
	if !config.CurrentConfig.Voices.Azure.Enabled {
		log.Print("WARN: azure is disabled")
		return
	}
	var err error
	request, err = http.NewRequest(http.MethodPost, baseUrl+"v1", nil)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Add("Ocp-Apim-Subscription-Key", config.CurrentConfig.Voices.Azure.Key)
	request.Header.Add("Content-Type", "application/ssml+xml")
	request.Header.Add("X-Microsoft-OutputFormat", "ogg-48khz-16bit-mono-opus")

	getReq, err := http.NewRequest(http.MethodGet, baseUrl+"voices/list", nil)
	if err != nil {
		log.Fatal(err)
	}
	getReq.Header.Add("Ocp-Apim-Subscription-Key", config.CurrentConfig.Voices.Azure.Key)
	cli := http.Client{}
	res, err := cli.Do(getReq)
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
	azureVoices = map[string]azureVoice{}
	for _, v := range voices {
		azureVoices[v.ShortName] = v
	}
}

func AzureSynth(content *string, voice *string) (*[]byte, error) {
	val, exists := azureVoices[*voice]
	if !exists {
		return nil, errors.New("invalid voice type")
	}
	buffer := bytes.Buffer{}
	err := xml.EscapeText(&buffer, []byte(*content))
	if err != nil {
		return nil, err
	}
	tmpReq := request
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
	return &bin, nil
}

func AzureVerify(voice *string) error {
	_, exists := azureVoices[*voice]
	if !exists {
		return errors.New("invalid voice type")
	}
	return nil
}
