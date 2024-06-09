package handler

import (
	"zeroChat/apps/task/mq/internal/handler/msgTransfer"
	"zeroChat/apps/task/mq/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

type Listen struct {
	svc *svc.ServiceContext
}

func NewListen(svc *svc.ServiceContext) *Listen {
	return &Listen{svc: svc}
}

// kafka 消息的消费者的添加
func (l *Listen) Services() []service.Service {
	return []service.Service{
		// 此处可以加载多个消费者
		kq.MustNewQueue(l.svc.Config.MsgReadTransfer, msgTransfer.NewMsgReadTransfer(l.svc)),
		kq.MustNewQueue(l.svc.Config.MsgChatTransfer, msgTransfer.NewMsgChatTransfer(l.svc)),
	}
}
