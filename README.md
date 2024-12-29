# Introduction 

GHIssueMarket sandbox, is controlled virtual environment for SWE-Agents' economic experimentation, simulating the environment of an envisioned peer-to-peer multi-agent system for GitHub issues outsourcing auctions. In this controlled setting, autonomous SWE-Agents auction and bid on GitHub issues, leveraging real-time communication, a built-in Retrieval-Augmented Generation (RAG) interface for effective decision-making, and instant cryptocurrency micropayments.

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


# Paper
[Intelligent Software Engineering Economics](https://arxiv.org/abs/2412.11722)
```
@misc{fouad2024ghissuemarketsandboxenvironmentsweagents,
      title={GHIssuemarket: A Sandbox Environment for SWE-Agents Economic Experimentation}, 
      author={Mohamed A. Fouad and Marcelo de Almeida Maia},
      year={2024},
      eprint={2412.11722},
      archivePrefix={arXiv},
      primaryClass={cs.SE},
      url={https://arxiv.org/abs/2412.11722}, 
}
```
