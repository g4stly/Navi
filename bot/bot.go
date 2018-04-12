package bot

import (
	"fmt"
	"time"
	"plugin"
	"github.com/bwmarrin/discordgo"
	"github.com/g4stly/navi/common"
)


type Bot struct {
	ID		string
	session		*discordgo.Session
	commands	map[string]Command
	modules		[]Module
	enabledModules	map[string]bool
	permissions	map[string]int
	Quit		chan int
}

type Command func(*Bot, *discordgo.User, int, []string) (string, error)
type Module struct {
	Name		string
	Commands	map[string]Command
}

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
	self.enabledModules[module.Name] = true

	return nil
}

func (self *Bot) loadCommands() {
	for _, module := range self.modules {
		common.Log("loading commands for module: %v", module.Name)
		if !self.enabledModules[module.Name] { continue }
		for name, command := range module.Commands {
			common.Log("\t\tloading %v...", name)
			self.commands[name] = command
		}
	}
}

func (self *Bot) ReloadModules() error {
	common.Log("attepmting to hot load modules")

	// reload config file
	err := common.LoadConfig()
	if err != nil { return err }

	// loop over names of modules
	modules := common.Config["modules"].([]interface{})
	for moduleIndex, moduleName := range modules {
		common.Log("loading module: %v", moduleName.(string))

		// check to see if it's already loaded
		if _, ok := self.enabledModules[moduleName.(string)]; ok {
			// if so, delete it first
			delete(self.enabledModules, moduleName.(string))
			self.modules = append(self.modules[:moduleIndex], self.modules[moduleIndex+1:]...)
		}

		err = self.loadModule(moduleName.(string))
		if err != nil { return err }
		common.Log("done")
	}
	// delete all loaded commands
	for name := range self.commands {
		delete(self.commands, name)
	}
	// finally reload commands
	self.loadCommands()
	return nil
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

func New(token string) (*Bot, error) {
	common.Log("Initializing bot...")
	var err error

	// initialize bot 
	bot := &Bot{
		commands:	make(map[string]Command),
		enabledModules:	make(map[string]bool),
		permissions:	make(map[string]int),
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
