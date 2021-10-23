package main

import (
	"Kotone-DiVE/lib"
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/db"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

func init() {
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
	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates
	err = discord.Open()
	if err != nil {
		log.Print("Discordgo connection failure:", err)
		return
	}

	log.Print("Kotone-DiVE started successfly!")
	discord.UpdateGameStatus(0, "Kotone-DiVE is running! | .help")
	defer discord.Close()
	defer db.Close()
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Print("Kotone-DiVE is gracefully shutdowning!")
}
