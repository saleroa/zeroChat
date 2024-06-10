package websocket

import (
	"math"
	"time"
)

const (
	// 心跳检测时间间隔
	defaultMaxConnectionIdle = time.Duration(math.MaxInt64)
	// ack 超时时间
	defaultAckTimeout  = 30 * time.Second
	defaultConcurrency = 10
)
