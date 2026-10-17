PLUGIN_DIR := plugins
BUILTIN_DIR := builtins

BIN_DIR := bin
SERVER := $(BIN_DIR)/server

PLUGIN_SRCS := $(wildcard $(BUILTIN_DIR)/*/)
PLUGINS := $(patsubst $(BUILTIN_DIR)/%/,$(PLUGIN_DIR)/%.so,$(PLUGIN_SRCS))

.PHONY: all clean run

all: $(SERVER) $(PLUGINS)

$(PLUGIN_DIR)/%.so: $(BUILTIN_DIR)/%/
	@mkdir -p $(PLUGIN_DIR)
	go build -buildmode=plugin -o $@ ./$<

$(SERVER): $(wildcard cmd/server/*.go)
	@mkdir -p $(BIN_DIR)
	go build -o $@ ./cmd/server

run: all
	./$(SERVER)

clean:
	rm -f $(PLUGIN_DIR)/*.so $(SERVER)
