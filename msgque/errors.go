package msgque

import "errors"

var (
	ErrMsgQSuspended     = errors.New("message queue suspended")
	ErrMsgQClosed        = errors.New("message queue closed")
	ErrMsgQEnqueOvertime = errors.New("message enqueue process overtime")
)

var (
	ErrCacheUpFailed = errors.New("cache up message failed")
)

var (
	ErrCallbackSendTimeout    = errors.New("message callback send result timeout")
	ErrCallbackReceiveTimeout = errors.New("message callback receive result timeout")
)
