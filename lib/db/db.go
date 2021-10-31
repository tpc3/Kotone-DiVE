package db

import (
	"Kotone-DiVE/lib/config"
	"errors"
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
	userCache       map[string]*config.User
)

func init() {
	var err error
	guildCache = map[string]*config.Guild{}
	userCache = map[string]*config.User{}
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

func LoadGuild(id *string) config.Guild {
	var (
		err   error
		guild *config.Guild
	)
	val, exists := guildCache[*id]
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
		guildCache[*id] = guild
		return *guild
	}
}

func SaveGuild(id *string, guild *config.Guild) error {
	var err error
	switch config.CurrentConfig.Db.Kind {
	case Bbolt:
		err = SaveGuildBbolt(id, guild)
	}
	if err != nil {
		log.Print("WARN: SaveGuild error:", err)
	} else {
		delete(guildCache, *id)
		delete(RegexCache, *id)
	}
	return err
}

func LoadUser(id *string) (config.User, error) {
	var (
		err  error
		user *config.User
	)
	val, exists := userCache[*id]
	if exists {
		return *val, nil
	} else {
		switch config.CurrentConfig.Db.Kind {
		case Bbolt:
			user, err = LoadUserBbolt(id)
		}
		if err != nil {
			log.Print("WARN: UserConfig is not available:", err)
			return config.User{}, err
		} else if user == nil {
			return config.User{}, errors.New("user does not exists")
		}
		userCache[*id] = user
		return *user, nil
	}
}

func SaveUser(id *string, user *config.User) error {
	var err error
	switch config.CurrentConfig.Db.Kind {
	case Bbolt:
		err = SaveUserBbolt(id, user)
	}
	if err != nil {
		log.Print("WARN: SaveUser error:", err.Error())
	} else {
		delete(userCache, *id)
	}
	return err
}

func DeleteUser(id *string) error {
	var err error
	switch config.CurrentConfig.Db.Kind {
	case Bbolt:
		err = DeleteUserBbolt(id)
	}
	if err != nil {
		log.Print("WARN: DeleteUser error:", err)
	}
	return err
}
