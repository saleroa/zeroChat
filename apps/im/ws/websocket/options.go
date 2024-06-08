package websocket

type ServerOptions func(opt *serverOption)

type serverOption struct {
	Authentication
	patten string
}

func newServerOption(opts ...ServerOptions) serverOption {
	o := serverOption{
		Authentication: new(authentication),
		patten:         "/ws",
	}

	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// 设置 auth
func WithServerAuthentication(auth Authentication) ServerOptions {
	return func(opt *serverOption) {
		opt.Authentication = auth
	}
}

// 设置 访问路径 patten
func WithServerPatten(patten string) ServerOptions {
	return func(opt *serverOption) {
		opt.patten = patten
	}
}
