package cmds

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/embed"
	"Kotone-DiVE/lib/voices"

	"github.com/bwmarrin/discordgo"
)

const Skip = "skip"

func SkipCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guildconf *config.Guild) {
	val, exists := db.StateCache[orgMsg.GuildID]
	if !exists {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guildconf.Lang, config.Lang[guildconf.Lang].Error.Joinfirst))
	} else {
		if val.Stream != nil {
			val.Stream.SetPaused(true)
			*val.Done <- voices.Skipped
			session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "⏩")
		}
	}
}
