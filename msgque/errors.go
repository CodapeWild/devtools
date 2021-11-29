package msgque

import "errors"

var (
	ErrMsgAssertionFailed = errors.New("message assertion failed")
	ErrTicketQueueNil     = errors.New("tickets queue can not be nil")
	ErrMsgQueSuspended    = errors.New("message queue suspended")
	ErrMsgQueClosed       = errors.New("message queue closed")
	ErrMsgEnqueOvertime   = errors.New("message enqueue process overtime")
	ErrMsgSendFailed      = errors.New("send message failed")
	ErrMsgSendFailedAll   = errors.New("send message failed after reach max retry times")
)

var (
	ErrCacheUpFailed = errors.New("cache up message failed")
)

var (
	ErrMsgRetryRunOut         = errors.New("message enqueue run out of retry times")
	ErrCallbackSendTimeout    = errors.New("message callback send result timeout")
	ErrCallbackReceiveTimeout = errors.New("message callback receive result timeout")
)
