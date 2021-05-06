package httpext

var (
	StateSuccess = &StdStatus{Status: 1000, Msg: "success"}
)

var (
	StateServiceAccessBlocked = &StdStatus{Status: 2000, Msg: "service access blocked"}
	StateParseParamFailed     = &StdStatus{Status: 2001, Msg: "parse parameter failed"}
	StateParamInvalid         = &StdStatus{Status: 2002, Msg: "invalid parameter for request"}
	StateProcessTimeout       = &StdStatus{Status: 2003, Msg: "processing timeout"}
	StateProcessError         = &StdStatus{Status: 2004, Msg: "processing error"}
	StateNotFound             = &StdStatus{Status: 2005, Msg: "not found"}
	StateVerifyFailed         = &StdStatus{Status: 2006, Msg: "data verification failed"}
	StateBackendAccessBlocked = &StdStatus{Status: 2007, Msg: "backend accessing blocked"}
	StateDataModifyForbidden  = &StdStatus{Status: 2008, Msg: "data modification forbidden"}
	StateTokenExpired         = &StdStatus{Status: 2009, Msg: "token expired"}
	StateDataSizeInvalid      = &StdStatus{Status: 2010, Msg: "data size invalid"}
	StateDataTypeInvalid      = &StdStatus{Status: 2011, Msg: "data type invalid"}
	StateDataMediaInvalid     = &StdStatus{Status: 2012, Msg: "data media invalid"}
)
