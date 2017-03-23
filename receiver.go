package easyhttp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"
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
	if err := u.req.ParseForm(); err != nil {
		return err
	}

	// 将receiver属性映射到map中
	fields := make(map[string]reflect.Value)
	v := reflect.ValueOf(u.receiver).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldInfo := v.Type().Field(i)
		tag := fieldInfo.Tag
		name := tag.Get("http")
		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		fields[name] = v.Field(i)
	}

	// 从req.Form取值更新上面map的值
	for name, values := range u.req.Form {
		f := fields[name]
		if !f.IsValid() {
			continue
		}

		if u.logger != nil {
			u.logger.Infof("%s: %v", name, values)
		}

		for _, value := range values {
			if f.Kind() == reflect.Slice {
				elem := reflect.New(f.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				f.Set(reflect.Append(f, elem))

			} else {
				if err := populate(f, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}

	// if stringSliceContent(u.req.Header["Content-type"], "application/json") {
	// }
	return u.unpackJSONParams()
}

func populate(v reflect.Value, value string) error {
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
	return nil
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
		u.logger.Infof("%v", body)
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
