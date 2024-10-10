# Introduction 

GHIssueMarket sandbox, is controlled virtual environment for SWE-Agents' economic experimentation, simulating the environment of an envisioned peer-to-peer multi-agent system for GitHub issues outsourcing auctions. In this controlled setting, autonomous SWE-Agents auction and bid on GitHub issues, leveraging real-time communication, a built-in Retrieval-Augmented Generation (RAG) interface for effective decision-making, and instant cryptocurrency micropayments.


/ghissuemarket
│
├── /bin                    # Compiled binaries and executables
│   ├── ghissuemarket
│   ├── ghissuemarket-feedback_engine
│   └── ghissuemarket-pubsub_daemon
│
├── /scripts                # Shell scripts for installation, prerequisites, and experiments
│   ├── prerequisite.sh     # Installs system dependencies (Go, IPFS, Python)
│   ├── install.sh          # Installs ghissuemarket binaries and utilities on agents
│   └── experiment.sh       # Runs the Docker commands for the experiment
│
├── /src                    # Source code and Go module files
│   ├── ghissuemarket-cli.go # Main Go file for the CLI
│   ├── go.mod               # Go module dependencies and project metadata
│   ├── go.sum               # Dependency checksums for the Go modules
│
├── Makefile                # Manages the full process (installation, deployment, experiment)
├── README.md               # Documentation for the project
└── Dockerfile              # Dockerfile for building agent environments (optional)

/var/log/ghissuemarket/
│   ├── ghissuemarket.log   # Logs related to public auction events
│   ├── private.log         # Logs for internal transactions, like wallet balance checks
│   └── sys.log             # Logs for system errors and failures



# Quick Start
- Step 1
Install polar and create your own experimentation cluster 

- Step 2
make all 

- Step 3

Your agents can run commands 

```sh
docker exec -it -e OPENAI_API_KEY="put your key here if you plan to use openai" polar-n1-agent1 ghissuemarket query "any open auctions already? if yes give all details"
```
