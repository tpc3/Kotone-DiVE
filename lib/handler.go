package lib

import (
	"Kotone-DiVE/lib/cmds"
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/embed"
	"Kotone-DiVE/lib/voices"
	"bytes"
	"io"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

func MessageCreate(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	guild := db.LoadGuild(orgMsg.GuildID)

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if orgMsg.Author.ID == session.State.User.ID {
		return
	}
	if strings.HasPrefix(orgMsg.Content, guild.Prefix) {
		switch strings.SplitN(orgMsg.Content, " ", 2)[0][1:] {
		case "ping":
			cmds.Ping(session, orgMsg)
		case "join":
			cmds.Join(session, orgMsg)
		case "leave":
			cmds.Leave(session, orgMsg)
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
	switch guild.Voice.Source {
	case "watson":
		bin, err = voices.Watson(&content, &guild.Voice.Type)
	default:
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, "Voice is not impremented:"+guild.Voice.Source))
		return
	}
	if err != nil {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, err))
		return
	}
	if bin == nil {
		//Nothing to read
		return
	}

	_, exists := db.VoiceLock[orgMsg.GuildID]
	if !exists {
		db.VoiceLock[orgMsg.GuildID] = &sync.Mutex{}
	}
	db.VoiceLock[orgMsg.GuildID].Lock()
	defer db.VoiceLock[orgMsg.GuildID].Unlock()
	db.ConnectionCache[orgMsg.GuildID].Speaking(true)
	defer db.ConnectionCache[orgMsg.GuildID].Speaking(false)
	//Send voice
	if config.CurrentConfig.Debug {
		log.Print(strconv.Itoa(len(*bin)), "bytes ogg audio.")
	}
	encode, err := dca.EncodeMem(bytes.NewReader(*bin), dca.StdEncodeOptions)
	defer encode.Cleanup()
	if err != nil {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, err))
	}
	done := make(chan error)
	dca.NewStream(encode, db.ConnectionCache[orgMsg.GuildID], done)
	err = <-done
	if err != nil && err != io.EOF {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, err))
	}
}
