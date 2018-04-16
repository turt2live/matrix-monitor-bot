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

# Prometheus Metrics

If metrics are enabled in your config, matrix-monitor-bot will serve up metrics for scraping by Prometheus. Every metric
that is exported is a [Histogram](https://prometheus.io/docs/concepts/metric_types/#histogram) metric. The following
metrics are exported:

* `monbot_ping_send_delay_seconds` - Number of seconds for the origin to send a ping to their homeserver
* `monbot_ping_receive_delay_seconds` - Number of seconds for a bot to receive a ping
* `monbot_ping_process_delay_seconds` - Number of seconds for a bot to process a ping event
* `monbot_pong_send_delay_seconds` - Number of seconds for the origin to send a pong in response to a ping to their homeserver
* `monbot_pong_receive_delay_seconds` - Number of seconds for a bot to receive a pong
* `monbot_ping_time_seconds` - Total number of seconds a ping lasts
* `monbot_pong_time_seconds` - Total number of seconds a pong lasts
* `monbot_rtt_seconds` - Total number of seconds for a given ping/pong sequence

### Example queries

**Average round trip time between two servers (t2bot.io -> matrix.org in this case)**
```
rate(monbot_rtt_seconds_sum{sourceDomain="t2bot.io",receivingDomain="matrix.org"}[2m]) / rate(monbot_rtt_seconds_count{sourceDomain="t2bot.io",receivingDomain="matrix.org"}[2m])
```

**Average time it takes a particular server to send a ping:**
```
rate(monbot_ping_send_delay_seconds_sum{sourceDomain="t2bot.io"}[2m]) / rate(monbot_ping_send_delay_seconds_count{sourceDomain="t2bot.io"}[2m])
```

**Average time it takes for a particular server to receive a ping:**
```
rate(monbot_ping_time_seconds_sum{receivingDomain="t2bot.io"}[2m]) / rate(monbot_ping_time_seconds_count{receivingDomain="t2bot.io"}[2m])
```


# Architecture

TODO: This section
* How the bot measures things
* What the Prometheus metrics are
* Why the bot uses m.room.message and not a custom event
* Why the bot uses messages for pongs instead of read receipts
* Why the display name gets overwritten and how it is used
* Why someone should run this on their server
