package db

import (
	"Kotone-DiVE/lib/config"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/patrickmn/go-cache"
)

var (
	guildCache      map[string]*config.Guild
	ConnectionCache map[string]*discordgo.VoiceConnection
	RegexCache      map[string]*map[*regexp.Regexp]*string
	VoiceCache      *cache.Cache
	VoiceLock       map[string]*sync.Mutex
	// userCache map[string]
)

func init() {
	var err error
	guildCache = map[string]*config.Guild{}
	ConnectionCache = map[string]*discordgo.VoiceConnection{}
	RegexCache = map[string]*map[*regexp.Regexp]*string{}
	VoiceCache = cache.New(24*time.Hour, 1*time.Hour)
	VoiceLock = map[string]*sync.Mutex{}
	switch config.CurrentConfig.Db.Kind {
	case Bbolt:
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
	case Bbolt:
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
		case Bbolt:
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
	case Bbolt:
		err = SaveGuildBbolt(id, guild)
	}
	if err != nil {
		log.Print("WARN: SaveGuild error:", err)
	} else {
		delete(guildCache, id)
		delete(RegexCache, id)
	}
	return err
}
