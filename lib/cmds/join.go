package cmds

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/embed"

	"github.com/bwmarrin/discordgo"
)

const Join = "join"

func JoinCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guildconf *config.Guild) {
	_, exists := db.ConnectionCache[orgMsg.GuildID]
	if exists {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guildconf.Lang, config.Lang[guildconf.Lang].Error.Join.Already))
	} else {
		guild, err := session.State.Guild(orgMsg.GuildID)
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guildconf.Lang, err))
		}
		for _, state := range guild.VoiceStates {
			if state.UserID == orgMsg.Author.ID {
				voice, err := session.ChannelVoiceJoin(orgMsg.GuildID, state.ChannelID, false, true)
				if err != nil {
					session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guildconf.Lang, config.Lang[guildconf.Lang].Error.Join.Failed))
				}
				db.ConnectionCache[orgMsg.GuildID] = voice
				session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "🖐")
				return
			}
		}
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guildconf.Lang, config.Lang[guildconf.Lang].Error.Join.Joinfirst))
	}
}
