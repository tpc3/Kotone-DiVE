package cmds

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
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
	req := strings.SplitN(orgMsg.Content, " ", 2)
	var (
		str string
		obj interface{}
	)
	if len(req) == 1 || req[1] != "user" {
		if len(req) == 1 || req[1] != "all" {
			g := guild
			g.PolicyList = nil
			g.Replace = nil
			obj = g
		} else {
			obj = guild
		}
	} else {
		var err error
		obj, err = db.LoadUser(&orgMsg.Author.ID)
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.SubCmd))
			return
		}
	}
	result, err := yaml.Marshal(obj)
	if err != nil {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
	}
	str = "```" + ext + "\n" + string(result) + "```"
	if len(str) > 2048 {
		session.ChannelFileSend(orgMsg.ChannelID, Dump+"."+ext, bytes.NewReader(result))
	} else {
		msg := embed.NewEmbed(session, orgMsg)
		msg.Title = strings.Title(Dump)
		msg.Description = str
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, msg)
	}
}
