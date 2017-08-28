package easyhttp

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// Mouthpiece 返回response的结果，记录错误日志
type Mouthpiece struct {
	resp    http.ResponseWriter
	Err     error  `json:"-"`
	Message string `json:"message"`
	Status  int    `json:"status"`

	Data interface{} `json:"data,omitempty"`
}

// NewMouthpiece 创建传话筒
func NewMouthpiece(resp http.ResponseWriter) (mp *Mouthpiece) {
	mp = new(Mouthpiece)
	mp.resp = resp
	mp.Status = -1
	return
}

// SetError 设置错误信息
func (mp *Mouthpiece) SetError(err error) {
	mp.Err = err
}

// Convey 将执行结果使用http response返回
func (mp *Mouthpiece) Convey() (err error) {
	if mp.Err != nil {
		mp.Status = -1
		mp.Message = mp.Err.Error()

	} else {
		mp.Status = 0
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

	respStr := string(respMsg)
	respStr, _ = strconv.Unquote(`"` + respStr + `"`)

	_, err = resp.Write([]byte(respStr))
	return
}
