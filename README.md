# TScli

*TScli* - a very simple terminal-based client for TSWeb online judge. It supports submitting problems and receiving
feedback on them.

# Installation

Use `go install github.com/khbminus/tscli` to install TScli into `$GOBIN` or `$GOHOME/bin` or `$HOME/go/bin`.

# Usage

```shell
# => login into TSWeb account and save credentials to ~/.tscli.json
tscli login
# Set contest for current local config (with dialog)
tscli local set-contest
# Set contest for current local config (with specific id)
tscli local set-contest --id="id"
# Create/Update local config(.tscli.local) and parse problems/compilers
tscli local parse
# Set compiler for submits  
tscli local set-compiler
# Show local config
tscli local show
# Submit file <ProblemId>.<ext> (e.g. 2A.cpp)
tscli submit <ProblemId>.<ext>
# Submit file and specify problem
tscli submit -p 2A file.extension
```

# Configuration

TScli uses two configuration files: global - `$HOME/.tscli.json` with credentials and `.tscli.local` (search in parent
directories is supported) with data specific to a contest. 