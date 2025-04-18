package gserr

const (
	ok                uint32 = 0
	unknownError      uint32 = 100000
	attachedMsgError  uint32 = 100001
	serverCommonError uint32 = 100002
	requestParamError uint32 = 100003
	tokenExpiredError uint32 = 100004
	noPermissionError uint32 = 100005
	dbNotRecordError  uint32 = 100006

	// user
	userExistError uint32 = 100011
)

var errMsgMap = map[uint32]string{
	ok:                "成功",
	unknownError:      "未知错误",
	serverCommonError: "系统繁忙，请稍后再试",
	requestParamError: "请求参数错误",
	tokenExpiredError: "登录已过期",
	noPermissionError: "没有权限",
	dbNotRecordError:  "没有记录",
	userExistError:    "用户已存在",
}

var (
	ErrUnknown      = New(unknownError, errMsgMap[unknownError])
	ErrServerCommon = New(serverCommonError, errMsgMap[serverCommonError])
	ErrRequestParam = New(requestParamError, errMsgMap[requestParamError])
	ErrTokenExpired = New(tokenExpiredError, errMsgMap[tokenExpiredError])
	ErrNoPermission = New(noPermissionError, errMsgMap[noPermissionError])
	ErrDBNotRecord  = New(dbNotRecordError, errMsgMap[dbNotRecordError])
	ErrUserExist    = New(userExistError, errMsgMap[userExistError])
)

func MsgFromCode(code uint32) (msg string, ok bool) {
	msg, ok = errMsgMap[code]
	return
}
