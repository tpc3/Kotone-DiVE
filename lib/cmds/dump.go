package cmds

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/embed"
	"bytes"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

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
	if len(req) == 1 {
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
		id := ""
		if req[1] == "user" {
			id = orgMsg.Author.ID
		} else if strings.HasPrefix(req[1], "<@") && len(orgMsg.Mentions) == 1 {
			id = orgMsg.Mentions[0].ID
		} else {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.SubCmd))
			return
		}
		obj, err = db.LoadUser(&id)
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
		msg.Title = cases.Title(language.Und, cases.NoLower).String(Dump)
		msg.Description = str
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, msg)
	}
}
