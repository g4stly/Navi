package bot

import (
	"fmt"
	"time"
	"plugin"
	"github.com/bwmarrin/discordgo"
	"github.com/g4stly/navi/common"
)


// types
type Bot struct {
	ID		string
	Combo		bool
	session		*discordgo.Session
	Commands	map[string]Command
	modules		[]Module
	EnabledModules	map[string]bool
	Permissions	map[string]bool
	Quit		chan int
}

//type Command func(*Bot, *discordgo.User, int, []string) (string, error)
type Command interface {
	Names() []string
	Usage(string) string
	Execute(*Bot, *discordgo.User, int, []string) (string, error)
}

type Module struct {
	Name		string
	Commands	[]Command
}

// variables
var ErrNotAuthorized = "**You aren't authorized to execute that command!**"

// methods for Bot type
func (self *Bot) loadModules() error {
	modules := common.Config["modules"].([]interface{})
	for _, moduleName := range modules {
		common.Log("loading module: %v", moduleName.(string))

		err := self.loadModule(moduleName.(string))
		if err != nil { return err }

		common.Log("done")
	}
	return nil
}

func (self *Bot) loadModule(moduleName string) error {
	moduleLocation := fmt.Sprintf("modules/%v/%v.so", moduleName, moduleName)
	moduleLibrary, err := plugin.Open(moduleLocation)
	if err != nil { return err }

	initialize, err := moduleLibrary.Lookup("Init")
	if err != nil { return err }

	module := initialize.(func()Module)()
	self.modules = append(self.modules, module)
	self.EnabledModules[module.Name] = true

	return nil
}

func (self *Bot) loadCommands() {
	for _, module := range self.modules {
		common.Log("loading commands for module: %v", module.Name)
		if !self.EnabledModules[module.Name] { continue }
		for _, command := range module.Commands {
			for _, name := range command.Names() {
				common.Log("\t\tloading %v...", name)
				self.Commands[name] = command
			}
		}
	}
}

func (self *Bot) Connect() int {
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

func (self *Bot) waitForQuit() int {
	for {
		select {
		case <-self.Quit:
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

// main exported function, initializes and returns a Bot type
func New(token string) (*Bot, error) {
	common.Log("Initializing bot...")
	var err error

	// initialize bot 
	bot := &Bot{
		Combo:		false,
		Commands:	make(map[string]Command),
		EnabledModules:	make(map[string]bool),
		Permissions:	make(map[string]bool),
		Quit:		make(chan int)}

	// load modules
	err = bot.loadModules()
	if err != nil {
		common.Fatal("failed to load module: %v", err)
	}
	bot.loadCommands()

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
