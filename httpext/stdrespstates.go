package httpext

type StdRespState struct {
	State int
	Msg   string
}

var (
	StateSuccess = &StdRespState{State: 20000, Msg: "process success"}
)

var (
	StateServiceAccessBlocked = &StdRespState{State: 10001, Msg: "service access blocked"}
	StateParseParamFailed     = &StdRespState{State: 10002, Msg: "parse parameters failed"}
	StateParamInvalid         = &StdRespState{State: 10003, Msg: "invalid parameter for request"}
	StateProcessTimeout       = &StdRespState{State: 10004, Msg: "processing timeout"}
	StateProcessError         = &StdRespState{State: 10005, Msg: "processing error"}
	StateNotFound             = &StdRespState{State: 10006, Msg: "not found"}
	StateVerifyFailed         = &StdRespState{State: 10007, Msg: "data verification failed"}
	StateBackendAccessBlocked = &StdRespState{State: 10008, Msg: "backend accessing blocked"}
	StateDataModifyForbidden  = &StdRespState{State: 10009, Msg: "data modification forbidden"}
	StateTokenExpired         = &StdRespState{State: 10010, Msg: "token expired"}
	StateDataSizeInvalid      = &StdRespState{State: 10012, Msg: "data size invalid"}
	StateDataTypeInvalid      = &StdRespState{State: 10013, Msg: "data type invalid"}
	StateDataMediaInvalid     = &StdRespState{State: 10014, Msg: "data media invalid"}
)
