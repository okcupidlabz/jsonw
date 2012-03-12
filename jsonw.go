// -*- mode: go; tab-width: 4; c-basic-offset: 4; indent-tabs-mode: nil; -*-

package jsonw

import (
    "fmt"
    "reflect"
)

type JsonWrap struct {
    dat interface{}
    err error
}

type JsonWrapError struct {
    msg string
}

func (e JsonWrapError) Error() string { return e.msg; }

func MakeJsonWrap (i interface{}) (jw *JsonWrap) {
    jw = new (JsonWrap);
    jw.dat = i;
    return jw;
}

func isInt(v reflect.Value) bool {
	k := v.Kind()
	return k == reflect.Int || k == reflect.Int8 || 
        k == reflect.Int16 || k == reflect.Int32 || 
        k == reflect.Int64
}

func isUint(v reflect.Value) bool {
	k := v.Kind()
	return k == reflect.Uint || k == reflect.Uint8 || 
        k == reflect.Uint16 || k == reflect.Uint32 || 
        k == reflect.Uint64
}

func (jw *JsonWrap) GetInt() (ret int64, err error) {
    v := reflect.ValueOf (jw.dat)
    if isInt (v) {
        ret = v.Int()
    } else if ! isUint (v) {
        err = JsonWrapError { "Field is not of integer type" }
    } else if v.Uint() <= (1<<63 - 1) {
        ret = int64(v.Uint());
    } else {
        err = JsonWrapError { "Signed int64 overflow error" }
    }
    return 
}

func (jw *JsonWrap) GetString() (string, bool) {
    s, ok := (jw.dat).(string);
    return s, ok;
}

func (jw *JsonWrap) AtIndex(i int) (*JsonWrap) {
    v, ok := (jw.dat).([]interface{});
    ret := new(JsonWrap);
    if ok && len(v) < i {
        ret.dat = v[i];
    }
    if !ok {
		ret.err = JsonWrapError { fmt.Sprintf("Array[%d] out of bounds", i) }
    }
    return ret;
}

func (jw *JsonWrap) IsNil() bool {
    return jw.dat == nil;
}

func (jw *JsonWrap) AtKey(s string) (*JsonWrap) {
    v, ok := (jw.dat).(map[string]interface{});
    ret := new (JsonWrap);
    if ok {
        val, ok := v[s];
        if ok {
            ret.dat = val;
        }
    }
    if !ok {
        ret.err = JsonWrapError { fmt.Sprintf ("Dict[%s] not found", s) }
    }
    return ret;
}
