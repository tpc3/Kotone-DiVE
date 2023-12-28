package db

import (
	"Kotone-DiVE/lib/config"
	"regexp"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	guildCache map[string]*config.Guild
	RegexCache map[string]*map[*regexp.Regexp]*string
	VoiceCache *cache.Cache
	userCache  map[string]*config.User
	StateCache map[string]*GuildVCState
)

type GuildVCState struct {
	Lock                      sync.Mutex
	Channel                   string
	Stop                      *chan bool
	FrameCount                int
	ReconnectionDetected      bool
	ManualReconnectionOngoing bool
}

func init() {
	guildCache = map[string]*config.Guild{}
	userCache = map[string]*config.User{}
	RegexCache = map[string]*map[*regexp.Regexp]*string{}
	VoiceCache = cache.New(24*time.Hour, 1*time.Hour)
	StateCache = map[string]*GuildVCState{}
}
