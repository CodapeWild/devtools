package msgque

import "errors"

var (
	ErrMsgQClosed        = errors.New("message queue closed")
	ErrMsgQEnqueOvertime = errors.New("message enqueue process overtime")
)

var (
	ErrCachePushFailed = errors.New("push message to cache failed")
)

var (
	ErrCallbackSendTimeout    = errors.New("message callback send result timeout")
	ErrCallbackReceiveTimeout = errors.New("message callback receive result timeout")
)
