package lib

import (
	"Kotone-DiVE/lib/cmds"
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/embed"
	"Kotone-DiVE/lib/voices"
	"io"
	"log"
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
		case cmds.Help:
			cmds.HelpCmd(session, orgMsg, &guild)
		case cmds.Policy:
			cmds.PolicyCmd(session, orgMsg, guild)
		case cmds.User:
			cmds.UserCmd(session, orgMsg, guild)
		}
		return
	}
	_, exists := db.ConnectionCache[orgMsg.GuildID]
	if exists {
		ttsHandler(session, orgMsg, &guild)
	}
}

func ttsHandler(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild) {
	if !guild.ReadBots && orgMsg.Author.Bot {
		return
	}
	content := orgMsg.Content
	if len(content) > guild.MaxChar {
		content = content[:guild.MaxChar]
	}
	if guild.ReadName {
		if orgMsg.Member.Nick != "" {
			content = orgMsg.Member.Nick + " " + content
		} else {
			content = strings.Split(orgMsg.Member.User.Username, "#")[0] + " " + content
		}
	}

	switch guild.Policy {
	case "allow":
		for k := range guild.PolicyList {
			if k == orgMsg.Author.ID {
				return
			}
		}
	case "deny":
		exists := false
		for k := range guild.PolicyList {
			if k == orgMsg.Author.ID {
				exists = true
			}
		}
		if !exists {
			return
		}
	}
	var voice *config.Voice
	user, err := db.LoadUser(orgMsg.Author.ID)
	if err != nil {
		voice = &guild.Voice
	} else {
		voice = &user.Voice
	}
	encoded, err := voices.GetVoice(session, voices.Replace(&orgMsg.GuildID, &guild.Replace, content), voice)
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

func VoiceStateUpdate(session *discordgo.Session, state *discordgo.VoiceStateUpdate) {
	alone := true
	guild, err := session.State.Guild(state.GuildID)
	if err != nil {
		log.Print("WARN: VoiceStateUpdate failed:", err)
	}
	for _, userState := range guild.VoiceStates {
		if state.ChannelID == userState.ChannelID && userState.UserID != session.State.User.ID {
			alone = false
		}
	}
	if alone {
		db.ConnectionCache[state.GuildID].Disconnect()
	}
}
