package httpext

import (
	"encoding/json"
	"net/http"
)

type StdStatus struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

var (
	StateSuccess = &StdStatus{Status: 0, Msg: "success"}
)

var (
	StateServiceAccessBlocked = &StdStatus{Status: 1, Msg: "service access blocked"}
	StateParseParamFailed     = &StdStatus{Status: 2, Msg: "parse parameter failed"}
	StateParamInvalid         = &StdStatus{Status: 3, Msg: "invalid parameter for request"}
	StateProcessTimeout       = &StdStatus{Status: 4, Msg: "processing timeout"}
	StateProcessError         = &StdStatus{Status: 5, Msg: "processing error"}
	StateDataNotFound         = &StdStatus{Status: 6, Msg: "data can not find"}
	StateDataVerifyFailed     = &StdStatus{Status: 7, Msg: "data verification failed"}
	StateDataAccessBlocked    = &StdStatus{Status: 8, Msg: "data access blocked"}
	StateDataModifyForbidden  = &StdStatus{Status: 9, Msg: "data modification forbidden"}
	StateDataExpired          = &StdStatus{Status: 10, Msg: "data expired"}
	StateDataSizeInvalid      = &StdStatus{Status: 11, Msg: "data size invalid"}
	StateDataTypeInvalid      = &StdStatus{Status: 12, Msg: "data type invalid"}
	StateDataMediaInvalid     = &StdStatus{Status: 13, Msg: "data media invalid"}
)

type StdResp struct {
	*StdStatus
	Payload interface{} `json:"payload"`
}

func NewStdResp(status *StdStatus, payload interface{}) *StdResp {
	return &StdResp{
		StdStatus: status,
		Payload:   payload,
	}
}

func (this *StdResp) WriteJson(respw http.ResponseWriter) (n int, err error) {
	var buf []byte
	if buf, err = json.Marshal(this); err != nil {
		respw.WriteHeader(http.StatusInternalServerError)
	} else {
		respw.Header().Set("Content-Type", "application/json")
		n, err = respw.Write(buf)
	}

	return
}
