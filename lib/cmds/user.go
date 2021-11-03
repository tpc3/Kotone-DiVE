package cmds

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/embed"
	"Kotone-DiVE/lib/voices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const User = "user"

func UserCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild config.Guild) {
	parsed := strings.SplitN(orgMsg.Content, " ", 3)
	user, err := db.LoadUser(&orgMsg.Author.ID)
	if err != nil {
		user = config.User{
			Voice: config.Voice{
				Source: "",
				Type:   "",
			},
			Name: "",
		}
	}
	if len(parsed) < 2 {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.SubCmd))
		return
	}
	switch parsed[1] {
	case "voice":
		if len(parsed) < 3 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.Value))
			return
		}
		options := strings.SplitN(parsed[2], " ", 2)
		if len(options) != 2 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.Value))
			return
		}
		err := voices.VerifyVoice(&options[0], &options[1], config.Lang[guild.Lang].Error.Voice)
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.Value+": "+err.Error()))
			return
		}
		user.Voice.Source = options[0]
		user.Voice.Type = options[1]
	case "name":
		if len(parsed) < 3 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.Value))
			return
		}
		user.Name = parsed[2]
	case "delete":
		err := db.DeleteUser(&orgMsg.Author.ID)
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
		}
		session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
		return

	default:
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.SubCmd))
		return
	}
	err = db.SaveUser(&orgMsg.Author.ID, &user)
	if err != nil {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
	}
	session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
}
