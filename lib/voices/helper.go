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
	configRegexBefore map[*regexp.Regexp]string
	configRegexAfter  map[*regexp.Regexp]string
	httpCli           *http.Client
	Skipped           error
	SourceDisabled    error
	NoSuchSource      error
)

type VoiceSource interface {
	Verify(voice string) error
	Synth(content string, voice *string) (*[]byte, error)
	GetInfo() VoiceInfo
}

type VoiceInfo struct {
	Type             string
	Format           string
	Container        string
	ReEncodeRequired bool
	Enabled          bool
}

func init() {
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
	SourceDisabled = errors.New("voice source is disabled")
	NoSuchSource = errors.New("no such voice source")
}

func SourceSwitcher(source *string) (VoiceSource, error) {
	log.Print(*source)
	var voiceSource VoiceSource
	switch *source {
	case Watson.Info.Type:
		voiceSource = Watson
	case Gtts.Info.Type:
		voiceSource = Gtts
	case Gcp.Info.Type:
		voiceSource = Gcp
	case Azure.Info.Type:
		voiceSource = Azure
	case VoiceText.Info.Type:
		voiceSource = VoiceText
	case Voicevox.Info.Type:
		voiceSource = Voicevox
	case Coeiroink.Info.Type:
		voiceSource = Coeiroink
	case AquestalkProxy.Info.Type:
		voiceSource = AquestalkProxy
	default:
		return nil, NoSuchSource
	}
	if !voiceSource.GetInfo().Enabled {
		return nil, SourceDisabled
	}
	return voiceSource, nil
}

func VerifyVoice(source *string, voice string) error {
	voiceSource, err := SourceSwitcher(source)
	if err != nil {
		return err
	}
	return voiceSource.Verify(voice)
}

func GetVoice(content *string, voice *config.Voice) (*dca.Decoder, error) {
	crc := strconv.FormatUint(crc64.Checksum([]byte(voice.Source+voice.Type+*content), crc64.MakeTable(crc64.ISO)), 10)
	_, exists := db.VoiceCache.Get(crc)
	if !exists {
		voiceSource, err := SourceSwitcher(&voice.Source)
		if err != nil {
			return nil, err
		}
		var bin *[]byte
		for i := 0; i < config.CurrentConfig.Voices.Retry; i++ {
			bin, err = voiceSource.Synth(*content, &voice.Type)
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

	err := session.VoiceConnections[orgMsg.GuildID].Speaking(true)
	if err != nil {
		return err
	}
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
		time.Sleep(100 * time.Millisecond) // Super-duper dirty hack
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
	Gcp.Close()
}
