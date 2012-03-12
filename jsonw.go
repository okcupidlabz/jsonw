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

func wrongType (w string, g reflect.Kind) JsonWrapError {
    return JsonWrapError { fmt.Sprintf("type error: wanted %s, got %s", w, g) }
}

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
    if jw.err != nil {
        err = jw.err;
    } else {
        v := reflect.ValueOf (jw.dat)
        if isInt (v) {
            ret = v.Int()
        } else if ! isUint (v) {
            err = wrongType ("int", v.Kind());
        } else if v.Uint() <= (1<<63 - 1) {
            ret = int64(v.Uint());
        } else {
            err = JsonWrapError { "Signed int64 overflow error" }
        }
    }
    return 
}

func (jw *JsonWrap) GetUint() (ret uint64, err error) {
    if jw.err != nil {
        err = jw.err;
    } else {
        v := reflect.ValueOf (jw.dat)
        if isUint (v) {
            ret = v.Uint()
        } else if ! isInt (v) {
            err = wrongType ("uint", v.Kind());
        } else if v.Int() >= 0 {
            ret = uint64(v.Int());
        } else {
            err = JsonWrapError { "Unsigned uint64 underflow error" }
        }
    }
    return
}

func (jw *JsonWrap) GetBool (ret bool, err error) {
    if jw.err != nil {
        err = jw.err
    } else {
        v := reflect.ValueOf (jw.dat)
        k := v.Kind()
        if k == reflect.Bool {
            ret = v.Bool();
        } else {
            err = wrongType("bool", k);
        }
    }
}

func (jw *JsonWrap) GetString() (ret string, err error) {
    if jw.err != nil {
        err = jw.err;
    } else {
        v := reflect.ValueOf (jw.dat)
        k := v.Kind()
        if k == reflect.String {
            ret = v.String();
        } else {
            err = wrongType("string", k);
        }
    }
    return
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
