// -*- mode: go; tab-width: 4; c-basic-offset: 4; indent-tabs-mode: nil; -*-

package jsonw

import (
    "fmt"
    "reflect"
)


type Wrapper struct {
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

func (i *Wrapper) getData() interface{} { return i.dat }
func (i *Wrapper) Error() *Error { return i.err; }
func (i *Wrapper) IsOk() bool { return i.Error() == nil; }

func NewWrapper (i interface{}) (rd *Wrapper) {
    rd = new (Wrapper);
    rd.dat = i;
    return rd;
}

func NewDictionary() *Wrapper {
    m := make (map[string]interface{})
    return NewWrapper (m);
}

func NewArray(l int) *Wrapper {
    m := make ([]interface{}, l)
    return NewWrapper (m);
}

func NewNil() *Wrapper {
    return NewWrapper(nil);
}

func NewInt(i int) *Wrapper {
    return NewWrapper(i);
}

func NewInt64(i int64) *Wrapper {
    return NewWrapper (i);
}

func NewUint64 (u uint64) *Wrapper {
    return NewWrapper(u)
}

func NewString (s string) *Wrapper {
    return NewWrapper(s);
}

func NewBool (b bool) *Wrapper {
    return NewWrapper(b);
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

func (rd *Wrapper) GetInt() (ret int64, err error) {
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

func (rd *Wrapper) GetUint() (ret uint64, err error) {
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

func (rd *Wrapper) GetBool() (ret bool, err error) {
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

func (rd *Wrapper) GetString() (ret string, err error) {
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

func (rd *Wrapper) AtIndex(i int) *Wrapper {
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

func (rd *Wrapper) Len() (ret int, err error) {
    tmp, v := rd.asArray()
    if v == nil {
        err = tmp.err
    } else {
        ret = len(v);
    }
    return
}

func (i *Wrapper) Keys() (v []string, err error) {
    tmp, d := i.asDictionary()
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

func (i *Wrapper) asArray() (ret *Wrapper, v []interface{}) {
    if i.err != nil {
        ret = i;
    } else {
        var ok bool
        v, ok = (i.dat).([]interface{});
        ret = new(Wrapper);
        if !ok {
            ret.err = wrongType ("array", reflect.ValueOf(i.dat).Kind());
        }
    }
    return
}

func (rd *Wrapper) IsNil() bool {
    return rd.dat == nil;
}

func (rd *Wrapper) AtKey(s string) *Wrapper {
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

func (w *Wrapper) SetKey(s string, val Wrapper) error {
    b, d := w.asDictionary()
    if d != nil {
        d[s] = val.getData()
    }
    return b.Error()
}

func (w *Wrapper) SetIndex (i int, val Wrapper) error {
    b, d := w.asArray()
    if d != nil {
        d[i] = val.getData()
    }
    return b.Error()

}

func (i *Wrapper) asDictionary() (ret *Wrapper, d map[string]interface{}) {
    if i.err != nil {
        ret = i
    } else {
        var ok bool
        d, ok = (i.dat).(map[string]interface{});
        ret = new (Wrapper);
        if !ok {
            ret.err = wrongType ("dict", reflect.ValueOf(i.dat).Kind());
        }
    }
    return 
}

