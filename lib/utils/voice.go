package utils

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Kotone-DiVE/lib/config"
	"github.com/tpc3/Kotone-DiVE/lib/db"
	"log"
	"time"
)

func VoiceDisconnect(session *discordgo.Session, guildID string) error {
	if db.StateCache[guildID].FrameCount != 0 {
		*db.StateCache[guildID].Stop <- true
		for i := 0; i < 5; i++ {
			if db.StateCache[guildID].FrameCount == 0 {
				break
			}
			time.Sleep(100 * time.Millisecond) // Super-duper dirty hack
		}
	}
	return session.GuildMemberMove(guildID, session.State.User.ID, nil)
}

func VoiceReconnect(session *discordgo.Session, guildID string, channelID string) {
	if db.StateCache[guildID].ManualReconnectionOngoing {
		return
	}
	db.StateCache[guildID].ManualReconnectionOngoing = true
	if config.CurrentConfig.Debug {
		log.Print("WARN: VoiceStateUpdate reconnecting to VC...")
	}
	for i := 0; i < config.CurrentConfig.Discord.Retry; i++ {
		_, err := session.ChannelVoiceJoin(guildID, channelID, false, true)
		if err != nil {
			log.Print("WARN: VoiceStateUpdate failed to join, retrying...:", err)
			session.GuildMemberMove(guildID, session.State.User.ID, nil)
		} else {
			break
		}
	}
	db.StateCache[guildID].ManualReconnectionOngoing = false
}
