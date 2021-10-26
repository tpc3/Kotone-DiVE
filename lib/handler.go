package lib

import (
	"Kotone-DiVE/lib/cmds"
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/embed"
	"Kotone-DiVE/lib/voices"
	"bytes"
	"hash/crc64"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/patrickmn/go-cache"
)

func init() {
	dca.Logger = nil
}

func MessageCreate(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	guild := db.LoadGuild(orgMsg.GuildID)

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if orgMsg.Author.ID == session.State.User.ID {
		return
	}
	if strings.HasPrefix(orgMsg.Content, guild.Prefix) {
		switch strings.SplitN(orgMsg.Content, " ", 2)[0][1:] {
		case cmds.Ping:
			cmds.PingCmd(session, orgMsg)
		case cmds.Join:
			cmds.JoinCmd(session, orgMsg, &guild)
		case cmds.Leave:
			cmds.LeaveCmd(session, orgMsg, &guild)
		case cmds.Dump:
			cmds.DumpCmd(session, orgMsg, &guild)
		case cmds.Config:
			cmds.ConfigCmd(session, orgMsg, guild)
		}
		return
	}
	_, exists := db.ConnectionCache[orgMsg.GuildID]
	if exists {
		ttsHandler(session, orgMsg, &guild)
	}
}

func ttsHandler(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild) {
	content := orgMsg.Content
	var (
		bin *[]byte
		err error
	)
	crc := strconv.FormatUint(crc64.Checksum([]byte(guild.Voice.Source+guild.Voice.Type+orgMsg.Content), crc64.MakeTable(crc64.ISO)), 10)
	val, exists := db.VoiceCache.Get(crc)
	if exists {
		bin = val.(*[]byte)
	} else {
		switch guild.Voice.Source {
		case voices.Watson:
			bin, err = voices.WatsonSynth(&content, &guild.Voice.Type)
		case voices.Gtts:
			bin, err = voices.GttsSynth(&content, &guild.Voice.Type)
		default:
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Guild.Voice))
			return
		}
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
			return
		}
		if bin == nil {
			//Nothing to read
			return
		} else {
			db.VoiceCache.Add(crc, bin, cache.DefaultExpiration)
		}
	}

	//Send voice
	if config.CurrentConfig.Debug {
		log.Print(strconv.Itoa(len(*bin)), " bytes audio.")
	}
	dca.Logger = nil
	encode, err := dca.EncodeMem(bytes.NewReader(*bin), dca.StdEncodeOptions)
	defer encode.Cleanup()
	if err != nil {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
		return
	}
	_, exists = db.VoiceLock[orgMsg.GuildID]
	if !exists {
		db.VoiceLock[orgMsg.GuildID] = &sync.Mutex{}
	}
	db.VoiceLock[orgMsg.GuildID].Lock()
	defer db.VoiceLock[orgMsg.GuildID].Unlock()
	db.ConnectionCache[orgMsg.GuildID].Speaking(true)
	defer db.ConnectionCache[orgMsg.GuildID].Speaking(false)
	done := make(chan error)
	dca.NewStream(encode, db.ConnectionCache[orgMsg.GuildID], done)
	err = <-done
	if err != nil && err != io.EOF {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
	}
}
