package config

import (
	"errors"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Debug   bool
	Replace map[string]string `yaml:",flow"`
	Help    string
	Discord struct {
		Token  string
		Status string
		Retry  int
	}
	Db struct {
		Kind string
		Path string
	}
	Voices struct {
		Retry  int
		Watson struct {
			Enabled bool
			Token   string
			Api     string
		}
		Gtts struct {
			Enabled bool
		}
		Gcp struct {
			Enabled bool
			Token   string
		}
		Azure struct {
			Enabled bool
			Key     string
			Region  string
		}
		VoiceText struct {
			Enabled bool
			Token   string
		}
		Voicevox struct {
			Enabled bool
			Api     string
		}
	}
	Guild Guild
	User  User
}

type Guild struct {
	Prefix       string
	Lang         string
	MaxChar      int
	Voice        Voice
	ReadBots     bool
	ReadName     bool
	ReadAllUsers bool
	Policy       string
	PolicyList   map[string]string
	Replace      map[string]string
}
type User struct {
	Voice Voice
	Name  string
}
type Voice struct {
	Source string
	Type   string
}

const configFile = "./config.yaml"

var CurrentConfig Config

func init() {
	loadLang()
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("Config load failed:", err)
	}
	err = yaml.Unmarshal(file, &CurrentConfig)
	if err != nil {
		log.Fatal("Config parse failed:", err)
	}

	//verify
	if CurrentConfig.Debug {
		log.Print("Debug is enabled")
	}
	if CurrentConfig.Voices.Retry < 1 {
		log.Fatal("Retry times must be 1+")
	}
	if CurrentConfig.Discord.Token == "" {
		log.Fatal("Token is empty")
	}
	err = VerifyGuild(&CurrentConfig.Guild)
	if err != nil {
		log.Fatal("Config verify failed:", err)
	}
}

//You should call voices.Verify before running this!
func VerifyGuild(guild *Guild) error {
	val, exists := Lang[guild.Lang]
	if !exists {
		return errors.New("no such language") //Don't use nil val!
	}
	guilderrorstr := val.Error.Guild
	if len(guild.Prefix) != 1 {
		return errors.New(guilderrorstr.Prefix)
	}
	if guild.MaxChar > 2000 {
		return errors.New(guilderrorstr.MaxChar)
	}
	if guild.Policy != "allow" && guild.Policy != "deny" {
		return errors.New(guilderrorstr.Policy)
	}
	return nil
}
