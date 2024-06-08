package websocket

// 消息结构体
type Message struct {
	Method string      `json:"method"`
	FromId string      `json:"fromid"`
	Data   interface{} `json:"data"`
}

func NewMessage(fromId string, data interface{}) *Message {
	return &Message{
		FromId: fromId,
		Data:   data,
	}
}
