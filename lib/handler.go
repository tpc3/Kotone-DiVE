package lib

import (
	"Kotone-DiVE/lib/cmds"
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/embed"
	"Kotone-DiVE/lib/voices"
	"io"
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
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
		case cmds.Replace:
			cmds.ReplaceCmd(session, orgMsg, guild)
		}
		return
	}
	_, exists := db.ConnectionCache[orgMsg.GuildID]
	if exists {
		ttsHandler(session, orgMsg, &guild)
	}
}

func ttsHandler(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild) {
	encoded, err := voices.GetVoice(session, voices.Replace(&orgMsg.GuildID, &guild.Replace, orgMsg.Content), &guild.Voice)
	if err != nil {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
	}

	_, exists := db.VoiceLock[orgMsg.GuildID]
	if !exists {
		db.VoiceLock[orgMsg.GuildID] = &sync.Mutex{}
	}
	db.VoiceLock[orgMsg.GuildID].Lock()
	defer db.VoiceLock[orgMsg.GuildID].Unlock()
	db.ConnectionCache[orgMsg.GuildID].Speaking(true)
	defer db.ConnectionCache[orgMsg.GuildID].Speaking(false)
	done := make(chan error)
	dca.NewStream(encoded, db.ConnectionCache[orgMsg.GuildID], done)
	err = <-done
	if err != nil && err != io.EOF {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
	}
}
