package gserr

const (
	ok                uint32 = 0
	PanicError        uint32 = 100000
	unknownError      uint32 = 100001
	attachedMsgError  uint32 = 100002
	ServerCommonError uint32 = 100003
	requestParamError uint32 = 100004
	tokenExpiredError uint32 = 100005
	noPermissionError uint32 = 100006
	dbNotRecordError  uint32 = 100007

	// user
	userExistError uint32 = 100101
	// file
	fileUploadRequiredError    uint32 = 100201
	fileMetaUninitializedError uint32 = 100202
	fileIncompleteChunkError   uint32 = 100203
)

var errMsgMap = map[uint32]string{
	ok:                         "成功",
	PanicError:                 "服务器内部错误",
	unknownError:               "未知错误",
	ServerCommonError:          "系统繁忙，请稍后再试",
	requestParamError:          "请求参数错误",
	tokenExpiredError:          "登录已过期",
	noPermissionError:          "没有权限",
	dbNotRecordError:           "没有记录",
	userExistError:             "用户已存在",
	fileUploadRequiredError:    "文件记录不存在，需要上传文件",
	fileMetaUninitializedError: "文件信息未初始化",
	fileIncompleteChunkError:   "分片不完整",
}

var (
	ErrUnknown               = New(unknownError, errMsgMap[unknownError])
	ErrAttachedMsgError      = New(attachedMsgError, errMsgMap[attachedMsgError])
	ErrServerCommon          = New(ServerCommonError, errMsgMap[ServerCommonError])
	ErrRequestParam          = New(requestParamError, errMsgMap[requestParamError])
	ErrTokenExpired          = New(tokenExpiredError, errMsgMap[tokenExpiredError])
	ErrNoPermission          = New(noPermissionError, errMsgMap[noPermissionError])
	ErrDBNotRecord           = New(dbNotRecordError, errMsgMap[dbNotRecordError])
	ErrUserExist             = New(userExistError, errMsgMap[userExistError])
	ErrFileUpload            = New(fileUploadRequiredError, errMsgMap[fileUploadRequiredError])
	ErrFileMetaUninitialized = New(fileMetaUninitializedError, errMsgMap[fileMetaUninitializedError])
	ErrFileIncompleteChunk   = New(fileIncompleteChunkError, errMsgMap[fileIncompleteChunkError])
)

func MsgFromCode(code uint32) (msg string, ok bool) {
	msg, ok = errMsgMap[code]
	return
}
