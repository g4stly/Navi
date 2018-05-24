TARGET 	= main
BOT	= bot/bot.go bot/callbacks.go bot/database.go
CORE	= modules/navi-core/navi-core.so
LUCK	= modules/navi-luck/navi-luck.so
REACT	= modules/navi-reaction/navi-reaction.so
MODULES	= $(CORE) $(LUCK) $(REACT)

all: $(TARGET) $(MODULES)

$(TARGET): common/common.go $(BOT)
	go build main.go
$(CORE): $(CORE:%.so=%.go) $(BOT)
	go build -buildmode=plugin -o $@ $(@:%.so=%.go)
$(LUCK): $(LUCK:%.so=%.go) $(BOT)
	go build -buildmode=plugin -o $@ $(@:%.so=%.go)
$(REACT): $(REACT:%.so=%.go) $(BOT)
	go build -buildmode=plugin -o $@ $(@:%.so=%.go)

.PHONY: clean test
clean: 
	rm -vf $(TARGET) $(MODULES)
