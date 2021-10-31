package cmds

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/embed"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const Policy = "policy"

func PolicyCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild config.Guild) {
	parsed := strings.SplitN(orgMsg.Content, " ", 3)
	if len(parsed) < 2 {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.SubCmd))
		return
	}
	switch parsed[1] {
	case "add":
		if len(orgMsg.Mentions) != 1 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.Value))
			return
		}
		guild.PolicyList[orgMsg.Author.ID] = orgMsg.Author.Username
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
		delete(guild.PolicyList, orgMsg.Author.ID)
	case "list":
		var keys []string
		for k := range guild.PolicyList {
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
			msg.Title = strings.Title(Replace)
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
		} else {
			session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
		}
	}
}
