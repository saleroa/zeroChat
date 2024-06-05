package rpcserver

import (
	"context"

	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zerr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func LogInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	res, err := handler(ctx, req)
	if err == nil {
		return res, nil
	}
	logx.WithContext(ctx).Errorf("【RPC SERVER ERR】 %v", err)
	// 透过层层封装找到最底层的 err
	causeErr := errors.Cause(err)
	// 接口断言，判断是不是 codemsg
	// 很多错误都是 xerr 中 new 出来的，new 返回的 error 实际上就是 codemsg 结构体
	if e, ok := causeErr.(*zerr.CodeMsg); ok {
		err = status.Error(codes.Code(e.Code), e.Msg)
	}

	return resp, err

}
