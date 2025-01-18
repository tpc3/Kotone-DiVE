package cmds

import (
	"context"
	"sync"
	"time"

	"github.com/tpc3/Kotone-DiVE/lib/config"
	"github.com/tpc3/Kotone-DiVE/lib/db"
	"github.com/tpc3/Kotone-DiVE/lib/embed"

	"github.com/bwmarrin/discordgo"
)

const Join = "join"

func JoinCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guildconf *config.Guild) {
	_, exists := db.StateCache[orgMsg.GuildID]
	if exists {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guildconf.Lang, config.Lang[guildconf.Lang].Error.Join.Already))
	} else {
		state, err := session.State.VoiceState(orgMsg.GuildID, orgMsg.Author.ID)
		if err == nil && state != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			voiceConn, err := session.ChannelVoiceJoin(ctx, orgMsg.GuildID, state.ChannelID, false, true)
			if err != nil {
				session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guildconf.Lang, config.Lang[guildconf.Lang].Error.Join.Failed))
			}
			if config.CurrentConfig.Debug {
				voiceConn.LogLevel = discordgo.LogDebug
			}
			db.StateCache[orgMsg.GuildID] = &db.GuildVCState{
				Lock:    sync.Mutex{},
				Channel: orgMsg.ChannelID,
			}
			session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "🖐")
			return
		} else {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guildconf.Lang, config.Lang[guildconf.Lang].Error.Joinfirst))
		}
	}
}
