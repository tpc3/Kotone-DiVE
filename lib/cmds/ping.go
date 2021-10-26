package cmds

import (
	"Kotone-DiVE/lib/embed"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const Ping = "ping"

func PingCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	msg := embed.NewEmbed(session, orgMsg)
	msg.Color = embed.ColorBlue
	msg.Title = strings.Title(Ping)
	msg.Description = "Pong!"
	session.ChannelMessageSendEmbed(orgMsg.ChannelID, msg)
}
