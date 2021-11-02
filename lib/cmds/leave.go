package cmds

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/embed"

	"github.com/bwmarrin/discordgo"
)

const Leave = "leave"

func LeaveCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild) {
	_, exists := db.StateCache[orgMsg.GuildID]
	if exists {
		err := db.StateCache[orgMsg.GuildID].Connection.Disconnect()
		delete(db.StateCache, orgMsg.GuildID)
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
			return
		}
		session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👋")
	} else {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Leave.None))
	}
}
