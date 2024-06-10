package constants

type MType int

const (
	TextMtype MType = iota
)

// 私发 群发
type ChatType int

const (
	GroupChatType ChatType = iota
	SingleChatType
)

type ContentType int

const (
	ContentChatMsg ContentType = iota
	ContentMakeRead
)
