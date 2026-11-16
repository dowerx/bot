PLUGIN_DIR := plugins
BUILTIN_DIR := builtins

BIN_DIR := bin
SERVER := $(BIN_DIR)/server

PLUGIN_SRCS := $(wildcard $(BUILTIN_DIR)/*/*.go)
PLUGIN_SHARED_SRCS := $(wildcard plugin/*.go)
SERVER_SRCS := $(wildcard cmd/server/*.go) $(wildcard bot/*.go) $(wildcard config/*.go) $(PLUGIN_SHARED_SRCS)

PLUGIN_DIRS := $(wildcard $(BUILTIN_DIR)/*/)
PLUGINS := $(patsubst $(BUILTIN_DIR)/%/,$(PLUGIN_DIR)/%.so,$(PLUGIN_DIRS))

.PHONY: all clean run

all: $(SERVER) $(PLUGINS)

$(PLUGIN_DIR)/%.so: $(PLUGIN_SRCS) $(PLUGIN_SHARED_SRCS) go.mod go.sum
	@mkdir -p $(PLUGIN_DIR)
	go build -buildmode=plugin -o $@ ./$(BUILTIN_DIR)/$*

$(SERVER): $(SERVER_SRCS) go.mod go.sum
	@mkdir -p $(BIN_DIR)
	go build -o $@ ./cmd/server

run: all
	./$(SERVER)

clean:
	rm -f $(PLUGIN_DIR)/*.so $(SERVER)
