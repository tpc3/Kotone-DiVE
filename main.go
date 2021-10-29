package main

import (
	"Kotone-DiVE/lib"
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"Kotone-DiVE/lib/voices"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/common-nighthawk/go-figure"
)

func init() {
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

	log.Print("Kotone-DiVE started successfly!")
	discord.UpdateGameStatus(0, config.CurrentConfig.Discord.Status)
	defer discord.Close()
	defer db.Close()
	defer voices.CleanVoice()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Print("Kotone-DiVE is gracefully shutdowning!")
}
