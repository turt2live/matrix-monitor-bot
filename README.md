# matrix-monitor-bot

[![#monitorbot:t2bot.io](https://img.shields.io/badge/matrix-%23monitorbot:t2bot.io-brightgreen.svg)](https://matrix.to/#/#monitorbot:t2bot.io)
[![TravisCI badge](https://travis-ci.org/turt2live/matrix-monitor-bot.svg?branch=master)](https://travis-ci.org/turt2live/matrix-monitor-bot)

A bot to measure latency between homeservers, as perceived by users.

# Installing

Assuming Go 1.9 is already installed on your PATH:
```bash
# Get it
git clone https://github.com/turt2live/matrix-monitor-bot
cd matrix-monitor-bot

# Set up the build tools
currentDir=$(pwd)
export GOPATH="$currentDir/vendor/src:$currentDir/vendor:$currentDir:"$GOPATH
go get github.com/constabulary/gb/...
export PATH=$PATH":$currentDir/vendor/bin:$currentDir/vendor/src/bin"

# Build it
gb vendor restore
gb build

# Configure it (edit monitor-bot.yaml to meet your needs)
cp config.sample.yaml monitor-bot.yaml

# Run it
bin/monitor_bot
```

### Installing in Alpine Linux

The steps are almost the same as above. The only difference is that `gb build` will not work, so instead use the following lines:
```bash
go build -o bin/monitor_bot ./src/github.com/turt2live/matrix-monitor-bot/cmd/monitor_bot/
```

# Docker

`/path/to/matrix-monitor-bot` should always be pointed to a folder that has your `monitor-bot.yaml` file in it. If the config
file does not exist, one will be created for you (and promptly not work because it doesn't have a valid config). A folder
named `logs` will also be created here (assuming you use the default configuration).


**From Docker Hub:**
```
docker run -p 8080:8080 -v /path/to/matrix-monitor-bot:/data turt2live/matrix-monitor-bot
```


**Build the image yourself:**
```
git clone https://github.com/turt2live/matrix-monitor-bot
cd matrix-monitor-bot
docker build -t matrix-monitor-bot .
docker run -p 8080:8080 -v /path/to/matrix-monitor-bot:/data matrix-monitor-bot
```

# Architecture

TODO: This section
* How the bot measures things
* What the Prometheus metrics are
* Why the bot uses m.room.message and not a custom event
* Why the bot uses messages for pongs instead of read receipts
* Why the display name gets overwritten and how it is used
* Why someone should run this on their server
