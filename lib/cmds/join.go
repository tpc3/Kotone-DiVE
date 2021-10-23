package cmds

import (
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/embed"

	"github.com/bwmarrin/discordgo"
)

func Join(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	_, exists := db.ConnectionCache[orgMsg.GuildID]
	if exists {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, "I'm already joined..."))
	} else {
		guild, err := session.State.Guild(orgMsg.GuildID)
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, err))
		}
		for _, state := range guild.VoiceStates {
			if state.UserID == orgMsg.Author.ID {
				voice, err := session.ChannelVoiceJoin(orgMsg.GuildID, state.ChannelID, false, true)
				if err != nil {
					session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, "Join failed! please check your guild permission!"))
				}
				db.ConnectionCache[orgMsg.GuildID] = voice
				return
			}
		}
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, "You have to join VC first."))
	}
}
