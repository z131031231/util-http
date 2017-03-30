package easyhttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
)

// NewUnpacker 创建request参数解析器
func NewUnpacker(
	req *http.Request, receiver interface{}, logger Logger) (
	unpacker *Unpacker) {
	unpacker = new(Unpacker)
	unpacker.req = req
	unpacker.receiver = receiver
	unpacker.logger = logger

	return
}

// Unpack 将request中的请求参数解析到结构体中
func (u *Unpacker) Unpack() (err error) {
	if err = u.unpackJSONParams(); err != nil {
		return
	}

	return u.unpackGetParams()
}

// unpackGetParams 解析GET参数到接收器中
func (u *Unpacker) unpackGetParams() (err error) {
	vars := mux.Vars(u.req)
	if err = u.req.ParseForm(); err != nil {
		return err
	}

	rt := reflect.TypeOf(u.receiver)
	rv := reflect.ValueOf(u.receiver)

	if rt.Kind() == reflect.Ptr && rt.Elem().Kind() == reflect.Struct {
		rt = rt.Elem()
		rv = rv.Elem()
	}

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		key := f.Tag.Get("json")
		if key == "" {
			key = f.Name
		}

		val := u.req.FormValue(key)
		if val == "" && vars != nil {
			val = vars[key]
		}

		if u.logger != nil {
			u.logger.Debugf("key:%v value:%v", key, val)
		}

		if err = populate(rv.Field(i), val); err != nil {
			return
		}
	}

	return
}

func populate(v reflect.Value, value string) (err error) {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)

	default:
		return fmt.Errorf("unsupported http value type: %s", v.Type())
	}

	return
}

func (u *Unpacker) unpackJSONParams() (err error) {
	if u.req == nil || u.req.Body == nil {
		return fmt.Errorf("request body 为空")
	}

	/* if u.req.Body != nil {
		defer u.req.Body.Close()
	} */

	body, err := ioutil.ReadAll(u.req.Body)
	if err != nil {
		return
	}

	if u.logger != nil {
		u.logger.Info(string(body))
	}

	if len(body) > 0 {
		return json.Unmarshal(body, u.receiver)
	}

	return
}

func stringSliceContent(strs []string, str string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}

	return false
}
