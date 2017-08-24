package easyhttp

import (
	"encoding/json"
	"net/http"
)

// Mouthpiece 返回response的结果，记录错误日志
type Mouthpiece struct {
	resp    http.ResponseWriter
	err     error
	Message string `json:"message"`
	Status  int    `json:"status"`

	Data interface{} `json:"data,omitempty"`
}

// Convey 将执行结果使用http response返回
func (mp *Mouthpiece) Convey() (err error) {
	if mp.err != nil {
		mp.Status = -1
		mp.Message = mp.err.Error()
	}

	err = Response(mp.resp, mp)
	return
}

// Response 将结果打包成json返回给http
func Response(resp http.ResponseWriter, result interface{}) (err error) {
	respMsg, err := json.Marshal(result)
	if err != nil {
		return
	}

	_, err = resp.Write(respMsg)
	return
}
