package config

import (
	"time"
)

const RemoteSendDelayThreshold = 1500 * time.Millisecond
const ReceiveDelayThreshold = 5 * time.Second
const RttWarningThreshold = 10 * time.Second
const ProcessingDelayThreshold = 10 * time.Millisecond // This is pretty relaxed
const PingInterval = 5 * time.Minute
const MissedPingTimeout = 12 * time.Minute
const PongTimeout = 1 * time.Minute
const PingTtl = 15 * time.Minute
const RealRttTolerance = 20 * time.Millisecond
const WebWarnStatusThreshold = 1500 * time.Millisecond
const WebAverageInterval = 15 * time.Minute