# Docker Agent Configuration
DOCKER_AGENTS := polar-n1-agent1 polar-n1-agent2 polar-n1-agent3 polar-n1-agent4

# Application details
BINARY_NAME = ghissuemarket
TOOLS = ghissuemarket-feedback_engine ghissuemarket-pubsub_daemon

# Directories
INSTALL_DIR = /usr/local/bin
LOG_DIR = /var/log/ghissuemarket
BIN_DIR = ./bin
SRC_DIR = ./src
SCRIPTS_DIR = /scripts

# Default target
all: install

# Prerequisite Setup on Docker Agents
prerequisite-setup:
	@echo "Copying prerequisites script to Docker agents..."
	@for agent in $(DOCKER_AGENTS); do \
		docker exec -it $$agent mkdir -p $(SCRIPTS_DIR); \
		docker cp ./scripts/prerequisite.sh $$agent:$(SCRIPTS_DIR)/prerequisite.sh; \
	done
	@echo "Running prerequisites setup on all Docker agents..."
	@for agent in $(DOCKER_AGENTS); do \
		echo "Installing prerequisites on $$agent..."; \
		docker exec -it $$agent bash -c "sh /scripts/prerequisite.sh"; \
	done
	@echo "Prerequisite setup completed on all agents."

# Build the ghissuemarket binary with Go modules
build:
	@echo "Initializing Go modules in $(SRC_DIR) directory..."
	cd $(SRC_DIR) && go mod init ghissuemarket || true
	@echo "Fetching required dependencies..."
	cd $(SRC_DIR) && go get github.com/spf13/cobra@latest
	cd $(SRC_DIR) && go get github.com/google/uuid@latest
	@echo "Building $(BINARY_NAME) binary..."
	cd $(SRC_DIR) && go build -o ../bin/$(BINARY_NAME) ghissuemarket-cli.go
	@echo "Build complete."

# Install ghissuemarket on Docker Agents
install: build
	@echo "Deploying ghissuemarket binary to agents..."
	@for agent in $(DOCKER_AGENTS); do \
		docker cp $(BIN_DIR)/$(BINARY_NAME) $$agent:$(INSTALL_DIR)/$(BINARY_NAME); \
		docker exec -it $$agent chmod +x $(INSTALL_DIR)/$(BINARY_NAME); \
	done
	@echo "Deploying pre-built utilities to agents..."
	@for agent in $(DOCKER_AGENTS); do \
		docker cp $(BIN_DIR)/ghissuemarket-feedback_engine $$agent:$(INSTALL_DIR)/ghissuemarket-feedback_engine; \
		docker cp $(BIN_DIR)/ghissuemarket-pubsub_daemon $$agent:$(INSTALL_DIR)/ghissuemarket-pubsub_daemon; \
		docker exec -it $$agent chmod +x $(INSTALL_DIR)/ghissuemarket-feedback_engine; \
		docker exec -it $$agent chmod +x $(INSTALL_DIR)/ghissuemarket-pubsub_daemon; \
	done
	@echo "Setting up log directories on all agents..."
	@for agent in $(DOCKER_AGENTS); do \
		docker exec -it $$agent bash -c "mkdir -p $(LOG_DIR) && touch $(LOG_DIR)/ghissuemarket.log $(LOG_DIR)/private.log $(LOG_DIR)/sys.log"; \
	done
	@$(MAKE) start-daemon
	@echo "ghissuemarket installation completed."

# Start the pubsub daemon on all agents
start-daemon:
	@echo "Starting ghissuemarket-pubsub_daemon on all Docker agents..."
	@for agent in $(DOCKER_AGENTS); do \
		docker exec -d $$agent /usr/local/bin/ghissuemarket-pubsub_daemon; \
	done
	@echo "Pubsub daemon started on all agents."

# Run the experiment
run-experiment:
	@echo "Running the experiment..."
	@bash ./scripts/experiment.sh $(DOCKER_AGENTS)
	@echo "Experiment completed."

# Clean up binaries and logs
clean:
	@echo "Cleaning up..."
	@rm -f $(BIN_DIR)/$(BINARY_NAME)
	@echo "Cleaned up."
