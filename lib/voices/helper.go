package voices

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"bytes"
	"errors"
	"hash/crc64"
	"io"
	"log"
	"regexp"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/patrickmn/go-cache"
)

var (
	Voices      []string
	configRegex map[*regexp.Regexp]string
)

func init() {
	Voices = []string{Watson, Gtts}
	configRegex = map[*regexp.Regexp]string{}
	for k, v := range config.CurrentConfig.Replace {
		configRegex[regexp.MustCompile(k)] = v
	}
}

func VerifyVoice(source *string, voice *string, voiceerror string) error {
	switch *source {
	case Watson:
		if !config.CurrentConfig.Voices.Watson.Enabled {
			return errors.New(voiceerror)
		}
		return WatsonVerify(voice)
	case Gtts:
		if !config.CurrentConfig.Voices.Gtts.Enabled {
			return errors.New(voiceerror)
		}
		return GttsVerify(voice)
	case Gcp:
		if !config.CurrentConfig.Voices.Gcp.Enabled {
			return errors.New(voiceerror)
		}
		return GcpVerify(voice)
	case Azure:
		if !config.CurrentConfig.Voices.Azure.Enabled {
			return errors.New(voiceerror)
		}
		return AzureVerify(voice)
	default:
		return errors.New(voiceerror)
	}
}

func GetVoice(session *discordgo.Session, message *string, voice *config.Voice) (*dca.Decoder, error) {
	var (
		err error
	)
	crc := strconv.FormatUint(crc64.Checksum([]byte(voice.Source+voice.Type+*message), crc64.MakeTable(crc64.ISO)), 10)
	_, exists := db.VoiceCache.Get(crc)
	if !exists {
		var bin *[]byte
		switch voice.Source {
		case Watson:
			if !config.CurrentConfig.Voices.Watson.Enabled {
				return nil, errors.New("voice is not available:" + Watson)
			}
			bin, err = WatsonSynth(message, &voice.Type)
		case Gtts:
			if !config.CurrentConfig.Voices.Gtts.Enabled {
				return nil, errors.New("voice is not available:" + Gtts)
			}
			bin, err = GttsSynth(message, &voice.Type)
		case Gcp:
			if !config.CurrentConfig.Voices.Gcp.Enabled {
				return nil, errors.New("voice is not available:" + Gcp)
			}
			bin, err = GcpSynth(message, &voice.Type)
		case Azure:
			if !config.CurrentConfig.Voices.Azure.Enabled {
				return nil, errors.New("voice is not available:" + Azure)
			}
			bin, err = AzureSynth(message, &voice.Type)
		default:
			return nil, errors.New("No such voice source:" + voice.Source)
		}
		if err != nil {
			return nil, err
		}
		if bin == nil {
			//Nothing to read
			return nil, nil
		} else {

			//Send voice
			if config.CurrentConfig.Debug {
				log.Print(strconv.Itoa(len(*bin)), " bytes audio.")
			}

			result, err := dca.EncodeMem(bytes.NewReader(*bin), dca.StdEncodeOptions)
			defer result.Cleanup()
			if err != nil {
				return nil, err
			} else {
				bin, err := io.ReadAll(result)
				if err != nil {
					return nil, err
				}
				db.VoiceCache.Add(crc, bin, cache.DefaultExpiration)
			}
		}
	}

	value, _ := db.VoiceCache.Get(crc)
	encoded := dca.NewDecoder(bytes.NewReader(value.([]byte)))
	return encoded, nil
}

func CleanVoice() {
	GcpClose()
}

func Replace(id *string, list *map[string]string, content string) *string {
	val, exists := db.RegexCache[*id]
	compiled := map[*regexp.Regexp]*string{}
	if exists {
		compiled = *val
	} else {
		for k, v := range *list {
			regex, err := regexp.Compile(k)
			if err == nil {
				compiled[regex] = &v
			}
		}
		db.RegexCache[*id] = &compiled
	}
	log.Print(compiled)
	for k, v := range compiled {
		content = k.ReplaceAllString(content, *v)
		if config.CurrentConfig.Debug {
			log.Print(content)
		}
	}
	for k, v := range configRegex {
		content = k.ReplaceAllString(content, v)
		if config.CurrentConfig.Debug {
			log.Print(content)
		}
	}
	return &content
}
