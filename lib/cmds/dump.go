package cmds

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/embed"
	"bytes"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v2"
)

const (
	Dump = "dump"
	ext  = "yaml"
)

func DumpCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild) {
	result, err := yaml.Marshal(guild)
	if err != nil {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
	}
	str := "```" + ext + "\n" + string(result) + "```"

	if len(str) > 2048 {
		session.ChannelFileSend(orgMsg.ChannelID, Dump+"."+ext, bytes.NewReader(result))
	}
	msg := embed.NewEmbed(session, orgMsg)
	msg.Title = strings.Title(Dump)
	msg.Description = str
	session.ChannelMessageSendEmbed(orgMsg.ChannelID, msg)
}
