package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/g4stly/navi/common"
	"time"
)

const commandPrefix = byte('.')

// actual callbacks
func (self *Bot) ready(s *discordgo.Session, r *discordgo.Ready) {
	common.Log("Got ready event!")
	self.ID = r.User.ID
}

func (self *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == self.ID && !self.Combo {
		return
	}
	self.parseMessage(m.Message)
}

// other bot methods
func (self *Bot) sendMessage(channelID string, message string) {
	self.session.ChannelTyping(channelID)
	time.Sleep(500 * time.Millisecond)
	self.session.ChannelMessageSend(channelID, message)
}

func (self *Bot) parseMessage(message *discordgo.Message) {
	msg := message.Content

	// one character messages can't be commands
	if len(msg) < 2 {
		return
	}

	// assert first character is the command prefix, and the second character is not
	if msg[0] != commandPrefix {
		return
	}
	if msg[1] == commandPrefix {
		return
	}

	var err error
	var response string
	argc, argv := parseArguments(msg[1:])

	// log command attempts just in case snoopy tries to hack me
	common.Log("%s (%s) said %s;\n\texecuting [%s]", message.Author.Username, message.Author.ID, msg, argv[0])

	// look up a command named the first word
	// in any message prefixed with our commandPrefix
	if cmd, ok := self.Commands[argv[0]]; ok {
		response, err = cmd.Execute(self, message.Author, argc, argv)
		if err != nil {
			common.Log("error during %v: %v", argv[0], err)
			response = "An error occured during the execution of that command. Please let bulb know."
		}
	} else { // was not found in our commands hashmap
		response = fmt.Sprintf("`%v` is not a valid command...", argv[0])
	}

	common.Log("\tresponse: %v", response)
	go self.sendMessage(message.ChannelID, response)
}

func (self *Bot) onShutdown() int {
	self.savePermissions()
	for index := range self.moduleCleanup {
		// call each modules destructor
		err := self.moduleCleanup[index](self)
		if err != nil {
			common.Log("moduleCleanup(): ", err)
		}
	}
	return 0
}

// second order functions
func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r'
}

func parseArguments(commandString string) (int, []string) {
	var argc int
	var argv []string
	length := len(commandString)

	for i := 0; i < length; i++ {
		c := commandString[i]
		if !isSpace(c) {
			var start, end int
			// if c is double quotes, and either the
			// first character of the message or is escaped
			// with a forwardslash
			if c == '"' && (i < 1 || commandString[i-1] != '\\') {
				i++
				start = i
				// while the current character isn't closing quotes
				// or escaped with a forwardslash increment i
				for i < length && (commandString[i] != '"' || commandString[i-1] == '\\') {
					i++
				}
				if commandString[i-1] == '\\' {
					end = i - 1
				} else {
					end = i
				}
			} else {
				start = i
				i++
				for i < length && !isSpace(commandString[i]) && (commandString[i] != '"' || commandString[i-1] == '\\') {
					i++
				}
				end = i
			}
			argv = append(argv, commandString[start:end])
			argc++
		}
	}

	return argc, argv
}
