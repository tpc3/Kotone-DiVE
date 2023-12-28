package cmds

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/embed"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	shellquote "github.com/kballard/go-shellquote"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const Replace = "replace"

func ReplaceCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild config.Guild) {
	parsed := strings.SplitN(orgMsg.Content, " ", 3)
	if len(parsed) < 2 {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.SubCmd))
		return
	}
	if guild.Replace == nil {
		guild.Replace = map[string]string{}
	}
	switch parsed[1] {
	case "set":
		if len(parsed) < 3 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Replace.Syntax))
			return
		}
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
		if len(parsed) < 3 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Replace.Syntax))
			return
		}
		_, exists := guild.Replace[parsed[2]]
		if exists {
			delete(guild.Replace, parsed[2])
		} else {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Replace.Del))
			return
		}
	case "delnum":
		if len(parsed) < 3 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Replace.Syntax))
			return
		}
		val, err := strconv.Atoi(parsed[2])
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Replace.Syntax))
			return
		} else if len(guild.Replace) <= val || val < 0 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Replace.Del))
			return
		}
		var keys []string
		for k := range guild.Replace {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		delete(guild.Replace, keys[val])

	case "list":
		if len(guild.Replace) == 0 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Replace.Empty))
			return
		}
		var keys []string
		for k := range guild.Replace {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		text := ""
		for i, v := range keys {
			text += "[" + strconv.Itoa(i) + "] \"" + v + "\" => \"" + guild.Replace[v] + "\"\n"
		}
		desc := "```\n" + text + "```"
		if len(desc) > 2048 {
			session.ChannelFileSend(orgMsg.ChannelID, "list.txt", strings.NewReader(text))
		} else {
			msg := embed.NewEmbed(session, orgMsg)
			msg.Title = cases.Title(language.Und, cases.NoLower).String(Replace)
			msg.Description = desc
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, msg)
		}
		return // readonly

	default:
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.SubCmd))
		return
	}
	err := db.SaveGuild(orgMsg.GuildID, &guild)
	if err != nil {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
	} else {
		session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
	}
}
