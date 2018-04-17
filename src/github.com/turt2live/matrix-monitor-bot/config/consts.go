package config

import (
	"time"
)

const RemoteSendDelayThreshold = 1 * time.Second
const ReceiveDelayThreshold = 5 * time.Second
const RttWarningThreshold = 10 * time.Second
const ProcessingDelayThreshold = 10 * time.Millisecond // This is pretty relaxed
const PingInterval = 1 * time.Minute
const MissedPingTimeout = 5 * time.Minute
const PongTimeout = 1 * time.Minute
const PingTtl = 10 * time.Minute
const RealRttTolerance = 10 * time.Millisecond
const WebWarnStatusThreshold = 1500 * time.Millisecond
