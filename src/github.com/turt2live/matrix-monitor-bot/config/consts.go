package config

import (
	"time"
)

const RemoteSendDelayWarnThreshold = 1500 * time.Millisecond
const ReceiveDelayWarnThreshold = 5 * time.Second
const PingTimeWarnThreshold = 1500 * time.Millisecond
const PingInterval = 5 * time.Minute
const WebWarnStatusThreshold = 1500 * time.Millisecond
const WebAverageInterval = 15 * time.Minute
