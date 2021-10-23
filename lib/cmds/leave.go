package cmds

import (
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/embed"

	"github.com/bwmarrin/discordgo"
)

func Leave(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	_, exists := db.ConnectionCache[orgMsg.GuildID]
	if exists {
		err := db.ConnectionCache[orgMsg.GuildID].Disconnect()
		delete(db.ConnectionCache, orgMsg.GuildID)
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, err))
		}
	} else {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, "No VC to leave"))
	}
}
