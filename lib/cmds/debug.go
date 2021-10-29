package cmds

import (
	"Kotone-DiVE/lib/config"
	"Kotone-DiVE/lib/voices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const Debug = "debug"

func DebugCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild) {
	parsed := strings.SplitN(orgMsg.Content, " ", 3)
	var str string
	if len(parsed) < 2 {
		str = ""
	} else {
		str = parsed[1]
	}
	_, trace := voices.Replace(&orgMsg.GuildID, &guild.Replace, str, true)
	session.ChannelFileSendWithMessage(orgMsg.ChannelID, "WARNING: Avoid using this future as much as you can.\n This will ruin the performance!", "debug.log", strings.NewReader(*trace))
}
