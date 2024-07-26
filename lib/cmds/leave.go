package cmds

import (
	"github.com/tpc3/Kotone-DiVE/lib/config"
	"github.com/tpc3/Kotone-DiVE/lib/db"
	"github.com/tpc3/Kotone-DiVE/lib/embed"
	"github.com/tpc3/Kotone-DiVE/lib/utils"
	"github.com/bwmarrin/discordgo"
)

const Leave = "leave"

func LeaveCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild) {
	_, exists := db.StateCache[orgMsg.GuildID]
	if exists {
		state, err := session.State.VoiceState(orgMsg.GuildID, session.State.User.ID)
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
			return
		}
		if state == nil || state.ChannelID == "" {
			// abnormal
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Leave.None))
			return
		}
		err = utils.VoiceDisconnect(session, orgMsg.GuildID)
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
			return
		}
		session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👋")
	} else {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Leave.None))
	}
}
