package lib

import (
	"Kotone-DiVE/lib/cmds"
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/embed"
	"Kotone-DiVE/lib/voices"
	"log"
	"runtime/debug"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

func init() {
	dca.Logger = nil
}

func MessageCreate(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	var start time.Time
	if config.CurrentConfig.Debug {
		start = time.Now()
	}

	guild := db.LoadGuild(&orgMsg.GuildID)
	defer func() {
		if err := recover(); err != nil {
			log.Print("Oops, ", err)
			debug.PrintStack()
		}
	}()

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if orgMsg.Author.ID == session.State.User.ID || orgMsg.Content == "" {
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
		case cmds.Debug:
			cmds.DebugCmd(session, orgMsg, &guild)
		case cmds.Skip:
			cmds.SkipCmd(session, orgMsg, &guild)
		}
		if config.CurrentConfig.Debug {
			log.Print("Processed in ", time.Since(start).Milliseconds(), "ms.")
		}
		return
	}
	_, exists := db.StateCache[orgMsg.GuildID]
	if exists {
		if db.StateCache[orgMsg.GuildID].Channel == orgMsg.ChannelID {
			ttsHandler(session, orgMsg, &guild)
		}
	}
	if config.CurrentConfig.Debug {
		log.Print("Processed in ", time.Since(start).Nanoseconds(), "ns.")
	}
}

func ttsHandler(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild) {

	if !guild.ReadBots && orgMsg.Author.Bot {
		return
	}

	content, err := orgMsg.ContentWithMoreMentionsReplaced(session)
	if err != nil {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
		return
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
	user, err := db.LoadUser(&orgMsg.Author.ID)

	if err != nil || user.Voice.Source == "" {
		voice = &guild.Voice
	} else {
		voice = &user.Voice
	}

	replaced, _ := voices.Replace(&orgMsg.GuildID, &guild.Replace, content, false)
	runeContent := []rune(*replaced)
	if len(runeContent) > guild.MaxChar {
		content = string(runeContent[:guild.MaxChar])
	} else {
		content = *replaced
	}

	db.StateCache[orgMsg.GuildID].Lock.Lock()
	defer db.StateCache[orgMsg.GuildID].Lock.Unlock()

	if guild.ReadName {
		var name *string
		if user.Name != "" {
			name = &user.Name
		} else if orgMsg.Member.Nick != "" {
			name = &orgMsg.Member.Nick
		} else {
			name = &strings.Split(orgMsg.Author.Username, "#")[0]
		}
		err = voices.ReadVoice(session, orgMsg, name, voice)
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
		}
	}
	err = voices.ReadVoice(session, orgMsg, &content, voice)
	if err != nil {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
	}
	db.StateCache[orgMsg.GuildID].Stream = nil
	db.StateCache[orgMsg.GuildID].Done = nil
	close(done)
}

func VoiceStateUpdate(session *discordgo.Session, state *discordgo.VoiceStateUpdate) {
	selfState, exists := db.StateCache[state.GuildID]
	if !exists {
		return // Bot isn't connected
	}
	alone := true
	guild, err := session.State.Guild(state.GuildID)
	if err != nil {
		log.Print("WARN: VoiceStateUpdate failed:", err)
	}
	for _, userState := range guild.VoiceStates {
		if selfState.Connection.ChannelID == userState.ChannelID && userState.UserID != session.State.User.ID {
			alone = false
		}
	}
	if alone {
		if db.StateCache[state.GuildID].Stream != nil {
			db.StateCache[state.GuildID].Stream.SetPaused(true)
			*db.StateCache[state.GuildID].Done <- io.EOF
			time.Sleep(100 * time.Millisecond) // Super duper dirty hack
		}
		err := db.StateCache[state.GuildID].Connection.Disconnect()
		delete(db.StateCache, state.GuildID)
		if err != nil {
			log.Print("WARN: VoiceStateUpdate failed to leave:", err)
		}
	}
}
