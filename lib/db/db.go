package db

import (
	"Kotone-DiVE/lib/config"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	guildCache      map[string]*config.Guild
	ConnectionCache map[string]*discordgo.VoiceConnection
	VoiceCache      map[string]map[string][]byte
	VoiceLock       map[string]*sync.Mutex
	// userCache map[string]
)

func init() {
	var err error
	guildCache = map[string]*config.Guild{}
	ConnectionCache = map[string]*discordgo.VoiceConnection{}
	VoiceCache = map[string]map[string][]byte{}
	VoiceLock = map[string]*sync.Mutex{}
	switch config.CurrentConfig.Db.Kind {
	case "bbolt":
		err = LoadBbolt()
	default:
		log.Fatal("That kind of db is not impremented:", config.CurrentConfig.Db.Kind)
	}
	if err != nil {
		log.Fatal("DB load error:", err)
	}
}

func Close() {
	var err error
	switch config.CurrentConfig.Db.Kind {
	case "bbolt":
		err = CloseBbolt()
	}
	if err != nil {
		log.Fatal("DB close error:", err)
	}
}
func LoadGuild(id string) config.Guild {
	var (
		err   error
		guild *config.Guild
	)
	val, exists := guildCache[id]
	if exists {
		return *val
	} else {
		switch config.CurrentConfig.Db.Kind {
		case "bbolt":
			guild, err = LoadGuildBbolt(id)
		}
		if err != nil {
			log.Print("WARN: LoadGuild error, using default:", err)
			return config.CurrentConfig.Guild
		}
		guildCache[id] = guild
		return *guild
	}
}

func SaveGuild(id string, guild config.Guild) error {
	var err error
	switch config.CurrentConfig.Db.Kind {
	case "bbolt":
		err = SaveGuildBbolt(id, guild)
	}
	if err != nil {
		log.Print("WARN: SaveGuild error:", err)
	}
	return err
}
