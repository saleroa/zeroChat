package websocket

import "time"

type ServerOptions func(opt *serverOption)

type serverOption struct {
	Authentication

	ack          AckType
	ackTimeout   time.Duration // ack 超时时间
	sendErrCount int           // ack 出错最大次数

	patten   string
	discover Discover

	maxConnectionIdle time.Duration

	concurrency int
}

func newServerOptions(opts ...ServerOptions) serverOption {
	o := serverOption{
		Authentication:    new(authentication),
		maxConnectionIdle: defaultMaxConnectionIdle,
		ackTimeout:        defaultAckTimeout,
		patten:            "/ws",
		concurrency:       defaultConcurrency,
		sendErrCount:      defaultSendErrCount,
	}

	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func WithServerAuthentication(auth Authentication) ServerOptions {
	return func(opt *serverOption) {
		opt.Authentication = auth
	}
}

func WithServerPatten(patten string) ServerOptions {
	return func(opt *serverOption) {
		opt.patten = patten
	}
}

func WithServerAck(ack AckType) ServerOptions {
	return func(opt *serverOption) {
		opt.ack = ack
	}
}

func WithServerMaxConnectionIdle(maxConnectionIdle time.Duration) ServerOptions {
	return func(opt *serverOption) {
		if maxConnectionIdle > 0 {
			opt.maxConnectionIdle = maxConnectionIdle
		}
	}
}

func WithSendErrCount(sendErrCount int) ServerOptions {
	return func(opt *serverOption) {
		opt.sendErrCount = sendErrCount
	}
}
func WithServerDiscover(discover Discover) ServerOptions {
	return func(opt *serverOption) {
		opt.discover = discover
	}
}
