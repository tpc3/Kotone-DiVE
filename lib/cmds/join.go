package cmds

import (
	"github.com/tpc3/Kotone-DiVE/lib/config"
	"github.com/tpc3/Kotone-DiVE/lib/db"
	"github.com/tpc3/Kotone-DiVE/lib/embed"
	"sync"

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
			_, err := session.ChannelVoiceJoin(orgMsg.GuildID, state.ChannelID, false, true)
			if err != nil {
				session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guildconf.Lang, config.Lang[guildconf.Lang].Error.Join.Failed))
			}
			db.StateCache[orgMsg.GuildID] = &db.GuildVCState{
				Lock:                      sync.Mutex{},
				Channel:                   orgMsg.ChannelID,
				ReconnectionDetected:      false,
				ManualReconnectionOngoing: false,
			}
			session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "🖐")
			return
		} else {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guildconf.Lang, config.Lang[guildconf.Lang].Error.Joinfirst))
		}
	}
}
