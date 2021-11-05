package voices

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"bytes"
	"errors"
	"hash/crc64"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/patrickmn/go-cache"
)

var (
	Voices      []string
	configRegex map[*regexp.Regexp]string
	httpCli     *http.Client
)

func init() {
	Voices = []string{Watson, Gtts}
	httpCli = &http.Client{}
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
	case VoiceText:
		if !config.CurrentConfig.Voices.VoiceText.Enabled {
			return errors.New(voiceerror)
		}
		return VoiceTextVerify(voice)
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
		case VoiceText:
			if !config.CurrentConfig.Voices.VoiceText.Enabled {
				return nil, errors.New("voice is not available:" + VoiceText)
			}
			bin, err = VoiceTextSynth(message, &voice.Type)
		default:
			return nil, errors.New("No such voice source:" + voice.Source)
		}
		if err != nil {
			return nil, err
		}
		if bin == nil {
			//Nothing to read
			return nil, nil
		}

		//Send voice
		if config.CurrentConfig.Debug {
			log.Print(strconv.Itoa(len(*bin)), " bytes audio.")
		}

		result, err := dca.EncodeMem(bytes.NewReader(*bin), dca.StdEncodeOptions)
		defer result.Cleanup()
		if err != nil {
			return nil, err
		}
		encBin, err := io.ReadAll(result)
		if err != nil {
			return nil, err
		}
		db.VoiceCache.Add(crc, encBin, cache.DefaultExpiration)
	}

	value, _ := db.VoiceCache.Get(crc)
	encoded := dca.NewDecoder(bytes.NewReader(value.([]byte)))
	return encoded, nil
}

func CleanVoice() {
	GcpClose()
}

func Replace(id *string, list *map[string]string, content string, trace bool) (*string, *string) {
	var (
		logStr string
		start  time.Time
		oldStr string
	) // debug
	if trace {
		start = time.Now()
		logStr = "Regex Replace() trace started at " + start.String() + " with string \"" + content + "\".\nGuildId is: " + *id + ".\n"
	}
	val, exists := db.RegexCache[*id]
	compiled := map[*regexp.Regexp]*string{}
	if exists {
		compiled = *val
		if trace {
			logStr += "Guild regex cache found.\nGot cache in " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns.\n"
		}
	} else {
		if trace {
			logStr += "Guild regex cache not found.\nCompiling regex...\n\n"
		}
		for k, v := range *list {
			if trace {
				logStr += "Compiling \"" + k + "\" ...\n"
			}
			regex, err := regexp.Compile(k)
			if err == nil {
				text := v //Let's encrypt knows everything
				compiled[regex] = &text
			} else {
				if trace {
					logStr += "|-Error occurred while compiling.\n" + err.Error() + ".\n|-Skipping...\n"
				}
			}
		}
		db.RegexCache[*id] = &compiled
		if trace {
			logStr += "Compiled regex in " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns.\n\n"
		}
	}
	if trace {
		logStr += "Starting process of " + strconv.Itoa(len(compiled)) + " user regex(s).\n\nRegex(s):\n"
		for k, v := range compiled {
			logStr += "|- \"" + k.String() + "\" => \"" + *v + "\"\n"
		}
		logStr += "\n"
	}
	for k, v := range compiled {
		if trace {
			oldStr = content
		}
		content = k.ReplaceAllString(content, *v)
		if trace && content != oldStr {
			logStr += "Regex hit!\n|-Regex: \"" + k.String() + "\"\n|-Replace: \"" + *v + "\"\n|-oldStr: \"" + oldStr + "\"\n|-New: \"" + content + "\"\n|-Time: " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns\n\n"
		}
	}
	if trace {
		logStr += "Processed user regex(s) in " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns.\n"
	}

	for k, v := range configRegex {
		if trace {
			oldStr = content
		}
		content = k.ReplaceAllString(content, v)
		if trace && content != oldStr {
			logStr += "Regex hit!\n|-Regex: \"" + k.String() + "\"\n|-Replace: \"" + v + "\"\n|-oldStr: \"" + oldStr + "\"\n|-New: " + content + "\n|-Time: " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns\n\n"
		}
	}
	if trace {
		logStr += "Processed config regex(s) in " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns.\nReplace() ended at " + time.Now().String() + " with string \"" + content + "\".\n"
	}
	return &content, &logStr
}
