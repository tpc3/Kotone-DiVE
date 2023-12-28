package embed

import (
	"Kotone-DiVE/lib/config"
	"log"
	"runtime/debug"
	"time"

	"github.com/bwmarrin/discordgo"
)

// https://material.io/archive/guidelines/style/color.html#color-color-palette
const (
	ColorBlue = 0x2196F3
	ColorPink = 0xf50057
)

var UnknownErrorNum int

func init() {
	UnknownErrorNum = 0
}

func NewEmbed(session *discordgo.Session, orgMsg *discordgo.MessageCreate) *discordgo.MessageEmbed {
	now := time.Now()
	msg := &discordgo.MessageEmbed{}
	msg.Author = &discordgo.MessageEmbedAuthor{}
	msg.Footer = &discordgo.MessageEmbedFooter{}
	msg.Author.IconURL = session.State.User.AvatarURL("256")
	msg.Author.Name = session.State.User.Username
	msg.Footer.IconURL = orgMsg.Author.AvatarURL("256")
	msg.Footer.Text = "Request from " + orgMsg.Author.Username + " @ " + now.String()
	msg.Color = ColorBlue
	return msg
}

func NewErrorEmbed(session *discordgo.Session, orgMsg *discordgo.MessageCreate, lang string, description string) *discordgo.MessageEmbed {
	msg := NewEmbed(session, orgMsg)
	msg.Color = ColorPink
	msg.Title = config.Lang[lang].Error.Title
	msg.Description = description
	return msg
}

func NewUnknownErrorEmbed(session *discordgo.Session, orgMsg *discordgo.MessageCreate, lang string, err error) *discordgo.MessageEmbed {
	log.Print("WARN: UnknownError called:", err)
	if config.CurrentConfig.Debug {
		debug.PrintStack()
	}
	UnknownErrorNum++
	return NewErrorEmbed(session, orgMsg, lang, config.Lang[lang].Error.Unknown)
}
