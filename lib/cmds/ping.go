package cmds

import (
	"Kotone-DiVE/lib/embed"

	"github.com/bwmarrin/discordgo"
)

func Ping(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	msg := embed.NewEmbed(session, orgMsg)
	msg.Color = embed.ColorBlue
	msg.Title = "Ping!"
	msg.Description = "Pong!"
	session.ChannelMessageSendEmbed(orgMsg.ChannelID, msg)
}
