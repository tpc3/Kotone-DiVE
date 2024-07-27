package cmds

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Kotone-DiVE/lib/config"
	"github.com/tpc3/Kotone-DiVE/lib/db"
	"github.com/tpc3/Kotone-DiVE/lib/embed"
)

const Skip = "skip"

func SkipCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guildconf *config.Guild) {
	val, exists := db.StateCache[orgMsg.GuildID]
	if !exists {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guildconf.Lang, config.Lang[guildconf.Lang].Error.Joinfirst))
	} else {
		if val.FrameCount != 0 {
			*val.Stop <- true
			session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "⏩")
		}
	}
}
