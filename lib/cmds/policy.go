package cmds

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/embed"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const Policy = "policy"

func PolicyCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild config.Guild) {
	parsed := strings.SplitN(orgMsg.Content, " ", 3)
	if len(parsed) < 2 {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.SubCmd))
		return
	}
	if guild.PolicyList == nil {
		guild.PolicyList = map[string]string{}
	}
	switch parsed[1] {
	case "add":
		if len(orgMsg.Mentions) != 1 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.Value))
			return
		}
		guild.PolicyList[orgMsg.Mentions[0].ID] = orgMsg.Mentions[0].Username
	case "del":
		if len(orgMsg.Mentions) != 1 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.Value))
			return
		}
		_, exists := guild.PolicyList[orgMsg.Mentions[0].ID]
		if !exists {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Policy.NotExists))
			return
		}
		delete(guild.PolicyList, orgMsg.Mentions[0].ID)
	case "list":
		if len(guild.Replace) == 0 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Replace.Empty))
			return
		}
		var keys []string
		for k := range guild.PolicyList {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		text := ""
		for i, v := range keys {
			text += "[" + strconv.Itoa(i) + "] \"" + v + "\" => \"" + guild.PolicyList[v] + "\"\n"
		}
		desc := "```\n" + text + "```"
		if len(desc) > 2048 {
			session.ChannelFileSend(orgMsg.ChannelID, "list.txt", strings.NewReader(text))
		} else {
			msg := embed.NewEmbed(session, orgMsg)
			msg.Title = cases.Title(language.Und, cases.NoLower).String(Policy)
			msg.Description = desc
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, msg)
		}
		return // readonly
	default:
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.SubCmd))
		return
	}

	err := config.VerifyGuild(&guild)
	if err != nil {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.Value+": "+err.Error()))
	} else {
		err = db.SaveGuild(&orgMsg.GuildID, &guild)
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
			return
		}
		session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")

	}
}
