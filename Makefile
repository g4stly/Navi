TARGET 	= main
BOT	= bot/bot.go bot/callbacks.go bot/database.go
CORE	= modules/navi-core/navi-core.so
MODULES	= $(CORE)

all: $(TARGET) $(MODULES)

$(TARGET): common/common.go $(BOT)
	go build main.go
$(CORE): $(CORE:%.so=%.go) $(BOT)
	go build -buildmode=plugin -o $@ $(@:%.so=%.go)

.PHONY: clean test
clean: 
	rm -vf $(TARGET) $(MODULES)
