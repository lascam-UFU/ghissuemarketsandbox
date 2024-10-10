#!/bin/bash

# Script to install ghissuemarket and its utilities on a Docker node
# This script assumes it is called with the agent name as an argument

set -e

AGENT="$1"
if [ -z "$AGENT" ]; then
    echo "Error: Agent name is not provided."
    exit 1
fi

# Paths and directories
INSTALL_DIR="/usr/local/bin"
LOG_DIR="/var/log/ghissuemarket"

# Build ghissuemarket binary
echo "Building ghissuemarket binary..."
go build -o ghissuemarket ./src/ghissuemarket-cli.go

# Copy binaries to the agent
echo "Deploying ghissuemarket binary and utilities to $AGENT..."
docker cp ghissuemarket "$AGENT:${INSTALL_DIR}/ghissuemarket"
docker cp ghissuemarket-feedback_engine "$AGENT:${INSTALL_DIR}/ghissuemarket-feedback_engine"
docker cp ghissuemarket-pubsub_daemon "$AGENT:${INSTALL_DIR}/ghissuemarket-pubsub_daemon"
docker exec -it $AGENT chmod +x ${INSTALL_DIR}/ghissuemarket
docker exec -it $AGENT chmod +x ${INSTALL_DIR}/ghissuemarket-feedback_engine
docker exec -it $AGENT chmod +x ${INSTALL_DIR}/ghissuemarket-pubsub_daemon

# Set up log directories on the agent
echo "Setting up log directories on $AGENT..."
docker exec -it $AGENT mkdir -p ${LOG_DIR}
docker exec -it $AGENT touch ${LOG_DIR}/ghissuemarket.log ${LOG_DIR}/private.log ${LOG_DIR}/sys.log
docker exec -it $AGENT chmod 755 ${LOG_DIR}
docker exec -it $AGENT chmod 664 ${LOG_DIR}/ghissuemarket.log ${LOG_DIR}/private.log ${LOG_DIR}/sys.log

echo "ghissuemarket installation completed successfully on $AGENT."
