package httpext

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"net/http"
)

type StdResp interface {
	Encode(state *StdRespState, payload interface{}) ([]byte, error)
	WriteBack(resp http.ResponseWriter, body []byte) (int, error)
}

type JsonResp struct{}

func (this *JsonResp) Encode(state *StdRespState, payload interface{}) ([]byte, error) {
	return json.Marshal(&struct {
		State   *StdRespState `json:"state"`
		Payload interface{}   `json:"payload"`
	}{
		State:   state,
		Payload: payload,
	})
}

func (this *JsonResp) WriteBack(resp http.ResponseWriter, body []byte) (int, error) {
	resp.Header().Set("Content-Type", "application/json")

	return resp.Write(body)
}

type GobResp struct{}

func (this *GobResp) Encode(state *StdRespState, payload interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(&struct {
		State   *StdRespState
		Payload interface{}
	}{
		State:   state,
		Payload: payload,
	})

	return buf.Bytes(), err
}

func (this *GobResp) WriteBack(resp http.ResponseWriter, body []byte) (int, error) {
	resp.Header().Set("Content-Type", "application/octet-stream")

	return resp.Write(body)
}
