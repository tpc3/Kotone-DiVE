package cmds

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/embed"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const Help = "help"

func HelpCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild) {
	msg := embed.NewEmbed(session, orgMsg)
	msg.Title = strings.Title(Help)
	msg.Description = config.Lang[guild.Lang].Help + "\n" + config.CurrentConfig.Help
	session.ChannelMessageSendEmbed(orgMsg.ChannelID, msg)
}
