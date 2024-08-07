package cmds

import (
	"github.com/tpc3/Kotone-DiVE/lib/db"
	"github.com/tpc3/Kotone-DiVE/lib/embed"
	"runtime"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const Ping = "ping"

func PingCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	msg := embed.NewEmbed(session, orgMsg)
	msg.Color = embed.ColorBlue
	msg.Title = cases.Title(language.Und, cases.NoLower).String(Ping)
	msg.Description = "Pong!"
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "Golang",
		Value: "`" + runtime.GOARCH + " " + runtime.GOOS + " " + runtime.Version() + "`",
	})
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "Stats",
		Value: "```\n" + strconv.Itoa(runtime.NumCPU()) + " cpu(s),\n" + strconv.Itoa(runtime.NumGoroutine()) + " go routine(s).```",
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "Memory",
		Value: "```\n" + strconv.FormatUint(mem.Sys/1024/1024, 10) + "MB used,\n" + strconv.FormatUint(uint64(mem.NumGC), 10) + " GCs.```",
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "Kotone-DiVE!",
		Value: "```\n" + strconv.Itoa(db.VoiceCache.ItemCount()) + " voices cached,\n" + strconv.Itoa(len(db.StateCache)) + " VCs ongoing,\n" + strconv.Itoa(embed.UnknownErrorNum) + " Unknown errors.```",
	})
	session.ChannelMessageSendEmbed(orgMsg.ChannelID, msg)
}
