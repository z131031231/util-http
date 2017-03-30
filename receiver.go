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
	if u.receiver == nil {
		return
	}

	if err = u.unpackJSONParams(); err != nil {
		return
	}

	return u.unpackGetParams()
}

// unpackGetParams 解析GET参数到接收器中
func (u *Unpacker) unpackGetParams() (err error) {
	if err = u.req.ParseForm(); err != nil {
		return err
	}

	rt := reflect.TypeOf(u.receiver)
	rv := reflect.ValueOf(u.receiver)

	if rt.Kind() == reflect.Ptr && rt.Elem().Kind() == reflect.Struct {
		return u.unpackFieldFromParams(rv)
	}

	return fmt.Errorf("解析参数类需要为 *struct 型，传入的是 %s", rt.String())

}

func (u *Unpacker) getFormVal(key string) (val string) {
	vars := mux.Vars(u.req)
	val = u.req.FormValue(key)
	if val == "" && vars != nil {
		val = vars[key]
	}

	return
}

func (u *Unpacker) unpackFieldFromParams(field reflect.Value) (err error) {
	rv := field.Elem()
	rt := field.Type().Elem()

	switch rt.Kind() {
	case reflect.Struct:
		for i := 0; i < rt.NumField(); i++ {
			f := rt.Field(i)
			key := f.Tag.Get("json")
			if key == "" {
				key = f.Name
			}

			val := u.getFormVal(key)
			if u.logger != nil {
				u.logger.Debugf("key:%v value:%v", key, val)
			}

			switch rv.Field(i).Kind() {
			case reflect.Ptr:
				u.unpackFieldFromParams(rv.Field(i).Addr())

			case reflect.Struct:
				u.unpackFieldFromParams(rv.Field(i).Addr())

			default:
				populate(rv.Field(i).Addr(), u.getFormVal(key))
			}

		}
	case reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			u.unpackFieldFromParams(rv.Index(i).Addr())
		}
	default:
		return fmt.Errorf("无法解析GET接收类型： %s", rt.String())
	}

	return
}

func populate(v reflect.Value, value string) (err error) {
	rv := v.Elem()
	switch v.Elem().Kind() {
	case reflect.String:
		rv.SetString(value)

	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		rv.SetInt(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		rv.SetBool(b)

	case reflect.Float32:
		f, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		rv.SetFloat(f)

	case reflect.Float64:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		rv.SetFloat(f)

	default:
		return fmt.Errorf("unsupported kind %s", v.Elem().Type().String())
	}

	return
}

func (u *Unpacker) unpackJSONParams() (err error) {
	if u.req == nil || u.req.Body == nil {
		return fmt.Errorf("request body 为空")
	}

	if u.req.Body != nil {
		defer u.req.Body.Close()
	}

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
