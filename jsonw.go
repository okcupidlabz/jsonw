// -*- mode: go; tab-width: 4; c-basic-offset: 4; indent-tabs-mode: nil; -*-

package jsonw

import (
    "fmt"
    "reflect"
)

type wrapper interface {
    IsOk() bool
    Error() *Error
    GetInt() (ret int64, err error)
    GetUint() (ret uint64, err error)
    GetBool() (ret bool, err error)
    GetString() (s string, err error)
    AtIndex(i int) wrapper
    Len() (ret int, err error) 
    Keys() (v []string, err error) 
    IsNil() bool
    AtKey(s string) wrapper
}

type Reader struct {
    dat interface{}
    err *Error
}

type Error struct {
    msg string
}

func (e Error) Error() string { return e.msg; }

func wrongType (w string, g reflect.Kind) *Error {
    return &Error { fmt.Sprintf("type error: wanted %s, got %s", w, g) }
}

func (rd *Reader) IsOk() bool { return rd.err == nil; }
func (rd *Reader) Error() *Error { return rd.err; }

func MakeReader (i interface{}) (rd *Reader) {
    rd = new (Reader);
    rd.dat = i;
    return rd;
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

func (rd *Reader) GetInt() (ret int64, err error) {
    if rd.err != nil {
        err = rd.err;
    } else {
        v := reflect.ValueOf (rd.dat)
        if isInt (v) {
            ret = v.Int()
        } else if ! isUint (v) {
            err = wrongType ("int", v.Kind());
        } else if v.Uint() <= (1<<63 - 1) {
            ret = int64(v.Uint());
        } else {
            err = Error { "Signed int64 overflow error" }
        }
    }
    return 
}

func (rd *Reader) GetUint() (ret uint64, err error) {
    if rd.err != nil {
        err = rd.err;
    } else {
        v := reflect.ValueOf (rd.dat)
        if isUint (v) {
            ret = v.Uint()
        } else if ! isInt (v) {
            err = wrongType ("uint", v.Kind());
        } else if v.Int() >= 0 {
            ret = uint64(v.Int());
        } else {
            err = Error { "Unsigned uint64 underflow error" }
        }
    }
    return
}

func (rd *Reader) GetBool() (ret bool, err error) {
    if rd.err != nil {
        err = rd.err
    } else {
        v := reflect.ValueOf (rd.dat)
        k := v.Kind()
        if k == reflect.Bool {
            ret = v.Bool();
        } else {
            err = wrongType("bool", k);
        }
    }
    return
}

func (rd *Reader) GetString() (ret string, err error) {
    if rd.err != nil {
        err = rd.err;
    } else {
        v := reflect.ValueOf (rd.dat)
        k := v.Kind()
        if k == reflect.String {
            ret = v.String();
        } else {
            err = wrongType("string", k);
        }
    }
    return
}

func (rd *Reader) AtIndex(i int) wrapper {
    ret, v := rd.asArray()
    if v == nil {

    } else if len (v) >= i {
        m := fmt.Sprintf ("index out of bounds %d >= %d", i, len(v))
        ret.err = &Error { m };
    } else {
        ret.dat = v[i];
    }
    return ret;
}

func (rd *Reader) Len() (ret int, err error) {
    rd, v := rd.asArray()
    if v == nil {
        err = rd.err
    } else {
        ret = len(v);
    }
    return
}

func (rd *Reader) Keys() (v []string, err error) {
    tmp, d := rd.asDictionary()
    if d == nil {
      err = tmp.err;
    } else {
      v = make([]string, len(d));
      var i int = 0;
      for k,_ := range d {
        v[i] = k
        i++
      }
    }
    return
}

func (rd *Reader) asArray() (ret *Reader, v []interface{}) {
    if rd.err != nil {
        ret = rd;
    } else {
        var ok bool
        v, ok = (rd.dat).([]interface{});
        ret = new(Reader);
        if !ok {
            ret.err = wrongType ("array", reflect.ValueOf(rd.dat).Kind());
        }
    }
    return
}

func (rd *Reader) IsNil() bool {
    return rd.dat == nil;
}

func (rd *Reader) AtKey(s string) wrapper {
    ret, d := rd.asDictionary()

    if d != nil {
        val,found := d[s];
        if found {
            ret.dat = val;
        } else {
            ret.dat = nil;
        }
    }
    return ret;
}

func (rd *Reader) asDictionary() (ret *Reader, d map[string]interface{}) {
    if rd.err != nil {
        ret = rd
    } else {
        var ok bool
        d, ok = (rd.dat).(map[string]interface{});
        ret = new (Reader);
        if !ok {
            ret.err = wrongType ("dict", reflect.ValueOf(rd.dat).Kind());
        }
    }
    return 
}

