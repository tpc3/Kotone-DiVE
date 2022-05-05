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
	parsed := strings.Split(orgMsg.Content, " ")
	p, err := session.State.MessagePermissions(orgMsg.Message)
	if err != nil {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
		return
	}
	admin := p&discordgo.PermissionAdministrator != 0
	if len(parsed) < 2 {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.SubCmd))
		return
	}
	id := &orgMsg.Author.ID
	if len(orgMsg.Mentions) != 0 {
		if !admin {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Permission))
			return
		}
		if !strings.HasPrefix(parsed[len(parsed)-1], "<@") {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.SubCmd))
			return
		}
		// change id, shrink parsed
		id = &orgMsg.Mentions[0].ID
		parsed = parsed[:len(parsed)-1]
	}
	user, err := db.LoadUser(id)
	if err != nil {
		user = config.User{
			Voice: config.Voice{
				Source: "",
				Type:   "",
			},
			Name: "",
		}
	}

	switch parsed[1] {
	case "voice":
		if len(parsed) != 4 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.Value))
			return
		}
		err := voices.VerifyVoice(&parsed[2], &parsed[3], config.Lang[guild.Lang].Error.Voice)
		if err != nil {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.Value+": "+err.Error()))
			return
		}
		user.Voice.Source = parsed[2]
		user.Voice.Type = parsed[3]
	case "name":
		if len(parsed) != 3 {
			session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewErrorEmbed(session, orgMsg, guild.Lang, config.Lang[guild.Lang].Error.Config.Value))
			return
		}
		user.Name = parsed[2]
	case "del":
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
	err = db.SaveUser(id, &user)
	if err != nil {
		session.ChannelMessageSendEmbed(orgMsg.ChannelID, embed.NewUnknownErrorEmbed(session, orgMsg, guild.Lang, err))
	}
	session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
}
