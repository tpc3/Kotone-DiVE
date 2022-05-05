package cmds

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/embed"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const Help = "help"

func HelpCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild) {
	msg := embed.NewEmbed(session, orgMsg)
	msg.Title = cases.Title(language.Und, cases.NoLower).String(Help)
	msg.Description = config.Lang[guild.Lang].Help + "\n" + config.CurrentConfig.Help
	session.ChannelMessageSendEmbed(orgMsg.ChannelID, msg)
}
