package cmds

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/embed"
	"Kotone-DiVE/lib/voices"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	Config = "config"
)

func ConfigCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild config.Guild) {
	parsed := strings.SplitN(orgMsg.Content, " ", 3)

	if len(parsed) < 2 {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.SubCmd))
		return
	}

	switch parsed[1] {
	case "prefix":
		if len(parsed) != 3 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.SubCmd))
			return
		}
		guild.Prefix = parsed[2]
	case "lang":
		if len(parsed) != 3 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.SubCmd))
			return
		}
		guild.Lang = parsed[2]
	case "maxchar":
		if len(parsed) != 3 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.SubCmd))
			return
		}
		i, err := strconv.Atoi(parsed[2])
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.Value))
			return
		} else {
			guild.MaxChar = i
		}
	case "voice":
		if len(parsed) != 3 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.SubCmd))
			return
		}
		opt := strings.SplitN(parsed[2], " ", 3)
		if len(opt) != 2 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.Value))
			return
		}
		err := voices.VerifyVoice(&opt[0], &opt[1], config.Lang[guild.Lang].Error.Guild.Voice)
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.Value+": "+err.Error()))
			return
		} else {
			guild.Voice.Source = opt[0]
			guild.Voice.Type = opt[1]
		}
	case "readbots":
		if len(parsed) != 3 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.SubCmd))
			return
		}
		i, err := strconv.ParseBool(parsed[2])
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.Value))
			return
		} else {
			guild.ReadBots = i
		}
	case "policy":
		if len(parsed) != 3 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.SubCmd))
			return
		}
		guild.Policy = parsed[2]
	default:
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.SubCmd))
		return
	}
	err := config.VerifyGuild(&guild)
	if err != nil {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.Value+": "+err.Error()))
	} else {
		err = db.SaveGuild(orgMsg.GuildID, guild)
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
			return
		} else {
			session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
		}
	}
}
