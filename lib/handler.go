package lib

import (
	"errors"
	"github.com/tpc3/Kotone-DiVE/lib/cmds"
	"github.com/tpc3/Kotone-DiVE/lib/config"
	"github.com/tpc3/Kotone-DiVE/lib/db"
	"github.com/tpc3/Kotone-DiVE/lib/embed"
	"github.com/tpc3/Kotone-DiVE/lib/utils"
	"github.com/tpc3/Kotone-DiVE/lib/voices"
	"log"
	"runtime/debug"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func MessageCreate(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	var start time.Time
	if config.CurrentConfig.Debug {
		start = time.Now()
	}

	guild := db.LoadGuild(orgMsg.GuildID)
	defer func() {
		if err := recover(); err != nil {
			log.Print("Oops, ", err)
			debug.PrintStack()
		}
	}()

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example, but it's a good practice.
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
		log.Print("Processed in ", time.Since(start).Milliseconds(), "ms.")
	}
}

func ttsHandler(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild) {

	if !guild.ReadBots && orgMsg.Author.Bot {
		return
	}

	if !guild.ReadAllUsers && !orgMsg.Author.Bot {
		state, err := session.State.VoiceState(orgMsg.GuildID, orgMsg.Author.ID)
		if errors.Is(err, discordgo.ErrStateNotFound) {
			return
		} else if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
			return
		}
		mystate, err := session.State.VoiceState(orgMsg.GuildID, session.State.User.ID)
		if errors.Is(err, discordgo.ErrStateNotFound) {
			return // ?????
		} else if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
			return
		}

		if state.ChannelID != mystate.ChannelID {
			return
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

	if err != nil || user.Voice.Source == "" {
		voice = &guild.Voice
	} else {
		voice = &user.Voice
	}

	content, err := orgMsg.ContentWithMoreMentionsReplaced(session)
	if err != nil {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
		return
	}

	replaced, _ := utils.Replace(orgMsg.GuildID, guild.Replace, content, false)
	if len(strings.TrimSpace(replaced)) == 0 {
		return
	}

	runeContent := []rune(replaced)
	if len(runeContent) > guild.MaxChar {
		content = string(runeContent[:guild.MaxChar])
	} else {
		content = replaced
	}

	var (
		encodedName    []byte
		encodedContent []byte
	)
	if guild.ReadName {
		var name *string
		if user.Name != "" {
			name = &user.Name
		} else if orgMsg.Member.Nick != "" {
			name = &orgMsg.Member.Nick
		} else {
			name = &strings.Split(orgMsg.Author.Username, "#")[0]
		}
		encodedName, err = voices.GetVoice(*name, voice)
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
		}
	}
	encodedContent, err = voices.GetVoice(content, voice)
	if err != nil {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
		return
	}
	if encodedContent == nil {
		return
	}

	db.StateCache[orgMsg.GuildID].Lock.Lock()
	defer func() {
		if db.StateCache[orgMsg.GuildID] != nil {
			db.StateCache[orgMsg.GuildID].Lock.Unlock()
		}
	}()

	if guild.ReadName && encodedName != nil {
		err = voices.ReadVoice(session, orgMsg, encodedName)
		if err != nil {
			if errors.Is(err, voices.Skipped) {
				return
			}
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
		}
	}
	err = voices.ReadVoice(session, orgMsg, encodedContent)
	if err != nil && !errors.Is(err, voices.Skipped) {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
	}
}

func VoiceStateUpdate(session *discordgo.Session, state *discordgo.VoiceStateUpdate) {
	if config.CurrentConfig.Debug {
		before := " ChannelBefore="
		if state.BeforeUpdate != nil {
			before += state.BeforeUpdate.ChannelID
		}
		log.Print("VoiceStateUpdate: UserID=" + state.UserID + before + " ChannelAfter=" + state.ChannelID + " SessionID=" + state.SessionID)
	}
	_, exists := db.StateCache[state.GuildID]
	if !exists {
		return // Bot isn't connected
	}
	if db.StateCache[state.GuildID].ManualReconnectionOngoing {
		return
	}
	_, exists = session.VoiceConnections[state.GuildID]

	guild, err := session.State.Guild(state.GuildID)
	if err != nil {
		log.Print("WARN: VoiceStateUpdate failed:", err)
	}

	myState, _ := session.State.VoiceState(state.GuildID, session.State.User.ID)
	if myState == nil {
		if db.StateCache[state.GuildID].ReconnectionDetected {
			log.Print("WARN: Will ignore this event due to detection")
			db.StateCache[state.GuildID].ReconnectionDetected = false
			return
		}
		if config.CurrentConfig.Debug {
			log.Print("I'm not exist in voicestates, maybe disconnection?")
		}
		delete(db.StateCache, state.GuildID)
		return
	}

	alone := true
	for _, userState := range guild.VoiceStates {
		if config.CurrentConfig.Debug {
			log.Print(userState)
		}
		if userState.UserID != session.State.User.ID {
			if myState.ChannelID == userState.ChannelID {
				alone = false
			}
		}
	}

	if alone {
		err = utils.VoiceDisconnect(session, guild.ID)
		if err != nil {
			log.Print("WARN: VoiceStateUpdate failed to leave:", err)
		}
		return
	}

	if state.UserID == session.State.User.ID && state.BeforeUpdate != nil {
		if state.BeforeUpdate.ChannelID != state.ChannelID && !exists {
			utils.VoiceReconnect(session, state.GuildID, state.ChannelID)
		} else if !db.StateCache[state.GuildID].ManualReconnectionOngoing {
			if state.BeforeUpdate.ChannelID == state.ChannelID && state.Suppress == state.BeforeUpdate.Suppress && state.SelfMute == state.BeforeUpdate.SelfMute && state.SelfDeaf == state.BeforeUpdate.SelfDeaf && state.Mute == state.BeforeUpdate.Mute && state.Deaf == state.BeforeUpdate.Deaf {
				log.Print("WARN: VoiceStateUpdate detected reconnection.")
				db.StateCache[state.GuildID].ReconnectionDetected = true
			}
		}
	}
}
