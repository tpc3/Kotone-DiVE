package main

import (
	"github.com/tpc3/Kotone-DiVE/lib"
	"github.com/tpc3/Kotone-DiVE/lib/config"
	"github.com/tpc3/Kotone-DiVE/lib/db"
	"github.com/tpc3/Kotone-DiVE/lib/voices"
	"log"
	"os"
	"os/signal"

	"net/http"
	_ "net/http/pprof"

	"github.com/bwmarrin/discordgo"
	"github.com/common-nighthawk/go-figure"
)

func init() {
	if config.CurrentConfig.Debug {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
	figure.NewColorFigure("Kotone-DiVE", "isometric1", "blue", true).Print()
	log.SetPrefix("[Init]")
	log.Print("Starting Kotone-DiVE!")
}

func main() {
	log.SetPrefix("[Main]")
	discord, err := discordgo.New("Bot " + config.CurrentConfig.Discord.Token)
	if err != nil {
		log.Fatal("Discordgo late init failure:", err)
	}
	discord.AddHandler(lib.MessageCreate)
	discord.AddHandler(lib.VoiceStateUpdate)
	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
	err = discord.Open()
	if err != nil {
		log.Print("Discordgo connection failure:", err)
		return
	}

	log.Print("Kotone-DiVE started successfully!")
	discord.UpdateGameStatus(0, config.CurrentConfig.Discord.Status)
	defer discord.Close()
	defer db.Close()
	defer voices.CleanVoice()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Print("Kotone-DiVE is gracefully shutdowning!")
}
