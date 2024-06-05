package xerr

var codeText = map[int]string{
	SERVER_COMMON_ERROR: "服务异常,请稍后处理",
	REQUEST_PARAM_ERROR: "参数错误",
	TOKEN_EXPIRE_ERROR:  "token失效,请重新登录",
	DB_ERROR:            "数据库繁忙,请稍后重试",
}

func ErrMsg(errcode int) string {
	if msg, ok := codeText[errcode]; ok {
		return msg
	}
	return codeText[SERVER_COMMON_ERROR]
}
