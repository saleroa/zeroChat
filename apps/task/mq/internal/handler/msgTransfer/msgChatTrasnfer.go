package msgTransfer

import (
	"context"
	"encoding/json"
	"fmt"
	"zeroChat/apps/im/immodels"
	"zeroChat/apps/im/ws/ws"
	"zeroChat/apps/task/mq/internal/svc"
	"zeroChat/apps/task/mq/mq"
	"zeroChat/pkg/bitmap"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// kafka 发送消息的消费者
type MsgChatTransfer struct {
	*baseMsgTransfer
}

func NewMsgChatTransfer(svc *svc.ServiceContext) *MsgChatTransfer {
	return &MsgChatTransfer{
		NewBaseMsgTransfer(svc),
	}
}

// kafka 消息的消费者
// 将消息加入到 mongo，然后 使用 push 方法推送到某个具体 conn 连接
// 实现了 consume 接口，作为 consumer
func (m *MsgChatTransfer) Consume(key, value string) error {
	fmt.Println("key : ", key, " value : ", value)

	var (
		data  mq.MsgChatTransfer
		ctx   = context.Background()
		msgId = primitive.NewObjectID()
	)
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}

	// 记录数据
	if err := m.addChatLog(ctx, msgId, &data); err != nil {
		return err
	}
	// 推送消息
	return m.Transfer(ctx, &ws.Push{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		RecvIds:        data.RecvIds,
		SendTime:       data.SendTime,
		MType:          data.MType,
		MsgId:          msgId.Hex(),
		Content:        data.Content,
	})
}

func (m *MsgChatTransfer) addChatLog(ctx context.Context, msgId primitive.ObjectID, data *mq.MsgChatTransfer) error {
	// 记录消息
	chatLog := immodels.ChatLog{
		ID:             msgId,
		ConversationId: data.ConversationId,
		SendId:         data.SendId,
		RecvId:         data.RecvId,
		ChatType:       data.ChatType,
		MsgFrom:        0,
		MsgType:        data.MType,
		MsgContent:     data.Content,
		SendTime:       data.SendTime,
	}

	readRecords := bitmap.NewBitmap(0)
	readRecords.Set(chatLog.SendId)
	chatLog.ReadRecords = readRecords.Export()

	err := m.svcCtx.ChatLogModel.Insert(ctx, &chatLog)
	if err != nil {
		return err
	}
	// 修改唯一的 conversation ，新增消息
	// 所有消息都会到这里，然后修改对应的 conversation
	return m.svcCtx.ConversationModel.UpdateMsg(ctx, &chatLog)
}
