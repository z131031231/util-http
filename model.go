package easyhttp

import "net/http"

// Logger 日志记录接口
type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
}

// Unpacker request参数解析器
type Unpacker struct {
	req      *http.Request
	receiver interface{}
	logger   Logger
}

type baseSender struct {
	url      string
	headers  map[string]string
	logger   Logger
	receiver interface{}
	rawResp  []byte
}

// GetSender get请求发送器
type GetSender struct {
	baseSender
	getParams map[string]string
}

// PostSender post请求发送器
type PostSender struct {
	GetSender
	postData interface{}
}
