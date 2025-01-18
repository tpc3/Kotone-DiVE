package utils

import (
	"github.com/bwmarrin/discordgo"
)

func VoiceDisconnect(session *discordgo.Session, guildID string) error {
	return session.VoiceStateUpdate(guildID, "", true, true)
}
