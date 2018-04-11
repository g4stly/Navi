package bot

import (
	"github.com/g4stly/navi/common"
	"github.com/bwmarrin/discordgo"
)

func (self *navi) ready(s *discordgo.Session, r *discordgo.Ready) {
	common.Log("Got ready event! Shutting down now")
	self.quit <- 1
}

func (self *navi) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// nothing yet
}

func (self *navi) onShutdown() int {
	common.Log("cleaing up in onShutdown()")
	return 0
}
