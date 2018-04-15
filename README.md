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
