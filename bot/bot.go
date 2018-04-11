package bot

import (
	"time"
	"github.com/bwmarrin/discordgo"
	"github.com/g4stly/navi/common"
)

type Bot interface {
	Connect() int
}

type command interface {
	Authorize(int) bool
	Execute(int, []string) (string, error)
}

// implements Bot
type navi struct {
	ID		string
	session		*discordgo.Session
	commands	map[string]command
	permissions	map[string]int
	quit		chan int
}

func (self *navi) Connect() int {
	common.Log("Connecting to Discord...")

	err := self.session.Open()
	if err != nil {
		common.Log("Failed to connect to Discord: %v", err)
		return -1
	}
	defer self.session.Close()

	common.Log("Done")
	return self.waitForQuit()
}

func (self *navi) waitForQuit() int {
	for {
		select {
		case <-self.quit:
			common.Log("Gracefully shutting down")
			resultCode := self.onShutdown()
			common.Log("Done")
			return resultCode;
		default:
			time.Sleep(500*time.Millisecond)
			break;
		}
	}
}

func New(token string) (Bot, error) {
	common.Log("Initializing bot...")
	var err error

	// initialize bot 
	bot := &navi{
		commands:	make(map[string]command),
		permissions:	make(map[string]int),
		quit:		make(chan int)}

	// create discord session
	bot.session, err = discordgo.New("Bot " + token)
	if err != nil {
		return nil, common.NewError("failed to create discordgo session: %v", err)
	}

	// add callbacks (defined in callbacks.go)
	bot.session.AddHandler(bot.ready)
	bot.session.AddHandler(bot.messageCreate)

	common.Log("Done")
	return bot, err;
}
