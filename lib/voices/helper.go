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
	Voices            []string
	configRegexBefore map[*regexp.Regexp]string
	configRegexAfter  map[*regexp.Regexp]string
	httpCli           *http.Client
	Skipped           error
)

func init() {
	Voices = []string{Watson, Gtts}
	httpCli = &http.Client{}
	configRegexBefore = map[*regexp.Regexp]string{}
	for k, v := range config.CurrentConfig.Replace.Before {
		configRegexBefore[regexp.MustCompile(k)] = v
	}
	configRegexAfter = map[*regexp.Regexp]string{}
	for k, v := range config.CurrentConfig.Replace.After {
		configRegexAfter[regexp.MustCompile(k)] = v
	}
	Skipped = errors.New("skipped")
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
	case Voicevox:
		if !config.CurrentConfig.Voices.Voicevox.Enabled {
			return errors.New(voiceerror)
		}
		return VoicevoxVerify(voice)
	case Coeiroink:
		if !config.CurrentConfig.Voices.Coeiroink.Enabled {
			return errors.New(voiceerror)
		}
		return CoeiroinkVerify(voice)
	case AquestalkProxy:
		if !config.CurrentConfig.Voices.AquestalkProxy.Enabled {
			return errors.New(voiceerror)
		}
		return AquestalkProxyVerify(voice)
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
		for i := 0; i < config.CurrentConfig.Voices.Retry; i++ {
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
			case Voicevox:
				if !config.CurrentConfig.Voices.Voicevox.Enabled {
					return nil, errors.New("voice is not available:" + Voicevox)
				}
				bin, err = VoicevoxSynth(message, &voice.Type)
			case Coeiroink:
				if !config.CurrentConfig.Voices.Coeiroink.Enabled {
					return nil, errors.New("voice is not available:" + Coeiroink)
				}
				bin, err = CoeiroinkSynth(message, &voice.Type)
			case AquestalkProxy:
				if !config.CurrentConfig.Voices.AquestalkProxy.Enabled {
					return nil, errors.New("voice is not available:" + AquestalkProxy)
				}
				bin, err = AquestalkProxySynth(message, &voice.Type)

			default:
				return nil, errors.New("No such voice source:" + voice.Source)
			}
			if err == nil {
				break
			}
		}
		if err != nil {
			return nil, err
		}

		if bin == nil {
			//Nothing to read
			return nil, nil
		}

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

func ReadVoice(session *discordgo.Session, orgMsg *discordgo.MessageCreate, encoded *dca.Decoder) error {
	if db.StateCache[orgMsg.GuildID].ManualReconnectionOngoing {
		for i := 0; i < config.CurrentConfig.Discord.Retry; i++ {
			if config.CurrentConfig.Debug {
				log.Print("Waiting for reconnection...")
			}
			time.Sleep(time.Second)
		}
	} else if session.VoiceConnections[orgMsg.GuildID] == nil {
		state, err := session.State.VoiceState(orgMsg.GuildID, session.State.User.ID)
		if err != nil {
			return err
		}
		VoiceReconnect(session, &orgMsg.GuildID, &state.ChannelID)
	}

	if session.VoiceConnections[orgMsg.GuildID] == nil {
		return Skipped //Skipped due to the disconnection
	}

	session.VoiceConnections[orgMsg.GuildID].Speaking(true)
	defer session.VoiceConnections[orgMsg.GuildID].Speaking(false)

	for i := 0; i < config.CurrentConfig.Discord.Retry; i++ {
		done := make(chan error)
		defer close(done)
		db.StateCache[orgMsg.GuildID].Done = &done
		db.StateCache[orgMsg.GuildID].Stream = dca.NewStream(encoded, session.VoiceConnections[orgMsg.GuildID], done)

		err := <-done

		_, exists := db.StateCache[orgMsg.GuildID]
		if !exists {
			return nil
		}
		db.StateCache[orgMsg.GuildID].Stream = nil
		db.StateCache[orgMsg.GuildID].Done = nil

		switch err {
		case io.EOF:
			fallthrough
		case nil:
			return nil
		case dca.ErrVoiceConnClosed:
			time.Sleep(time.Second)
			continue
		default:
			return err
		}
	}
	return dca.ErrVoiceConnClosed
}

func VoiceDisconnect(session *discordgo.Session, guildID *string) error {
	if db.StateCache[*guildID].Stream != nil {
		db.StateCache[*guildID].Stream.SetPaused(true)
		*db.StateCache[*guildID].Done <- io.EOF
		time.Sleep(100 * time.Millisecond) // Super duper dirty hack
	}
	return session.GuildMemberMove(*guildID, session.State.User.ID, nil)
}

func VoiceReconnect(session *discordgo.Session, guildID *string, channelID *string) {
	if db.StateCache[*guildID].ManualReconnectionOngoing {
		return
	}
	db.StateCache[*guildID].ManualReconnectionOngoing = true
	if config.CurrentConfig.Debug {
		log.Print("WARN: VoiceStateUpdate reconnecting to VC...")
	}
	for i := 0; i < config.CurrentConfig.Discord.Retry; i++ {
		_, err := session.ChannelVoiceJoin(*guildID, *channelID, false, true)
		if err != nil {
			log.Print("WARN: VoiceStateUpdate failed to join, retrying...:", err)
			session.GuildMemberMove(*guildID, session.State.User.ID, nil)
		} else {
			break
		}
	}
	db.StateCache[*guildID].ManualReconnectionOngoing = false
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

	for k, v := range configRegexBefore {
		if trace {
			oldStr = content
		}
		content = k.ReplaceAllString(content, v)
		if trace && content != oldStr {
			logStr += "Regex hit!\n|-Regex: \"" + k.String() + "\"\n|-Replace: \"" + v + "\"\n|-oldStr: \"" + oldStr + "\"\n|-New: " + content + "\n|-Time: " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns\n\n"
		}
	}
	if trace {
		logStr += "Processed config before regex(s) in " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns.\n"
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

	for k, v := range configRegexAfter {
		if trace {
			oldStr = content
		}
		content = k.ReplaceAllString(content, v)
		if trace && content != oldStr {
			logStr += "Regex hit!\n|-Regex: \"" + k.String() + "\"\n|-Replace: \"" + v + "\"\n|-oldStr: \"" + oldStr + "\"\n|-New: " + content + "\n|-Time: " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns\n\n"
		}
	}
	if trace {
		logStr += "Processed config after regex(s) in " + strconv.FormatInt(time.Since(start).Nanoseconds(), 10) + "ns.\nReplace() ended at " + time.Now().String() + " with string \"" + content + "\".\n"
	}
	return &content, &logStr
}
