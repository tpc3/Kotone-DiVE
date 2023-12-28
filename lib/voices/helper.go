package voices

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/utils"
	"errors"
	"hash/crc64"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/patrickmn/go-cache"
)

var (
	httpCli         *http.Client
	Skipped         error
	SourceDisabled  error
	NoSuchSource    error
	VoiceConnClosed error
)

type VoiceSource interface {
	Verify(voice string) error
	Synth(content string, voice string) ([]byte, error)
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
	Skipped = errors.New("skipped")
	SourceDisabled = errors.New("voice source is disabled")
	NoSuchSource = errors.New("no such voice source")
	VoiceConnClosed = errors.New("voice connection is closed")
}

func SourceSwitcher(source string) (VoiceSource, error) {
	var voiceSource VoiceSource
	switch source {
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

func VerifyVoice(source string, voice string) error {
	voiceSource, err := SourceSwitcher(source)
	if err != nil {
		return err
	}
	return voiceSource.Verify(voice)
}

func GetVoice(content string, voice *config.Voice) ([]byte, error) {
	crc := strconv.FormatUint(crc64.Checksum([]byte(voice.Source+voice.Type+content), crc64.MakeTable(crc64.ISO)), 10)
	_, exists := db.VoiceCache.Get(crc)
	if !exists {
		voiceSource, err := SourceSwitcher(voice.Source)
		if err != nil {
			return nil, err
		}
		var bin []byte
		for i := 0; i < config.CurrentConfig.Voices.Retry; i++ {
			bin, err = voiceSource.Synth(content, voice.Type)
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
			log.Print(strconv.Itoa(len(bin)), " bytes audio.")
		}

		result, err := Encode(bin, voiceSource.GetInfo())
		db.VoiceCache.Add(crc, result, cache.DefaultExpiration)
	}

	value, _ := db.VoiceCache.Get(crc)
	encoded := value.([]byte)
	return encoded, nil
}

func ReadVoice(session *discordgo.Session, orgMsg *discordgo.MessageCreate, encoded []byte) error {
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
		utils.VoiceReconnect(session, orgMsg.GuildID, state.ChannelID)
	}

	if session.VoiceConnections[orgMsg.GuildID] == nil {
		return Skipped //Skipped due to the disconnection
	}

	frames, err := SplitToFrame(encoded)
	if err != nil {
		return err
	}

	err = session.VoiceConnections[orgMsg.GuildID].Speaking(true)
	if err != nil {
		return err
	}
	defer session.VoiceConnections[orgMsg.GuildID].Speaking(false)

	for i := 0; i < config.CurrentConfig.Discord.Retry; i++ {
		stop := make(chan bool)
		defer close(stop)
		db.StateCache[orgMsg.GuildID].Stop = &stop
	FrameLoop:
		for i, v := range frames[db.StateCache[orgMsg.GuildID].FrameCount:] {
			select {
			case <-stop:
				err = Skipped
				break FrameLoop
			default:
				db.StateCache[orgMsg.GuildID].FrameCount = i
				if session.VoiceConnections[orgMsg.GuildID].Ready == false || session.VoiceConnections[orgMsg.GuildID].OpusSend == nil {
					err = VoiceConnClosed
					break FrameLoop
				}
				session.VoiceConnections[orgMsg.GuildID].OpusSend <- v
			}
		}

		_, exists := db.StateCache[orgMsg.GuildID]
		if !exists {
			return nil
		}
		db.StateCache[orgMsg.GuildID].Stop = nil

		switch {
		case err == io.EOF:
			fallthrough
		case err == nil:
			db.StateCache[orgMsg.GuildID].FrameCount = 0
			return nil
		case errors.Is(err, VoiceConnClosed):
			time.Sleep(time.Second)
			continue
		default:
			return err
		}
	}
	db.StateCache[orgMsg.GuildID].FrameCount = 0
	return VoiceConnClosed
}

func CleanVoice() {
	Gcp.Close()
}
