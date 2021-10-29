package cmds

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/embed"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	shellquote "github.com/kballard/go-shellquote"
)

const Replace = "replace"

func ReplaceCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild config.Guild) {
	parsed := strings.SplitN(orgMsg.Content, " ", 3)
	if len(parsed) < 2 {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.SubCmd))
		return
	}
	switch parsed[1] {
	case "set":
		options, err := shellquote.Split(parsed[2])
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Replace.Syntax))
			return
		}
		if len(options) != 2 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Replace.Syntax))
			return
		}
		_, err = regexp.Compile(options[0])
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Replace.Regex))
			return
		}
		guild.Replace[options[0]] = options[1]
	case "del":
		_, exists := guild.Replace[parsed[2]]
		if exists {
			delete(guild.Replace, parsed[2])
		} else {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Replace.Regex))
			return
		}
	case "delnum":
		val, err := strconv.Atoi(parsed[2])
		if err != nil || len(guild.Replace) < val || val < 0 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Replace.Syntax))
			return
		}
		var keys []string
		for k := range guild.Replace {
			keys = append(keys, k)
		}
		delete(guild.Replace, keys[val])

	case "list":
		var keys []string
		for k := range guild.Replace {
			keys = append(keys, k)
		}
		text := ""
		for i, v := range keys {
			text += "[" + strconv.Itoa(i) + "] \"" + v + "\" => \"" + guild.Replace[v] + "\"\n"
		}
		desc := "```\n" + text + "```"
		if len(desc) > 2048 {
			session.ChannelFileSend(orgMsg.ChannelID, "list.txt", strings.NewReader(text))
		} else {
			msg := embed.NewEmbed(session, orgMsg)
			msg.Title = strings.Title(Replace)
			msg.Description = desc
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, msg)
		}
		return // readonly

	default:
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.SubCmd))
		return
	}
	err := db.SaveGuild(orgMsg.GuildID, guild)
	if err != nil {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
	} else {
		session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
	}
}
