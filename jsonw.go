package jsonw

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Wrapper struct {
	dat    interface{}
	err    *Error
	access []string
}

type Error struct {
	msg string
}

func (e Error) Error() string { return e.msg }

func (w *Wrapper) NewError(format string, a ...interface{}) *Error {
	m1 := fmt.Sprintf(format, a...)
	p := w.AccessPath()
	m2 := fmt.Sprintf("%s: %s", p, m1)
	return &Error{m2}
}

func (w *Wrapper) wrongType(want string, got reflect.Kind) *Error {
	return w.NewError("type error: wanted %s, got %s", want, got)
}

func (i *Wrapper) getData() interface{} { return i.dat }
func (i *Wrapper) IsOk() bool           { return i.Error() == nil }

func (i *Wrapper) GetData() (dat interface{}, err error) {
	if i.err != nil {
		err = *i.err
	} else {
		dat = i.dat
	}
	return
}

func (i *Wrapper) GetDataVoid(dp *interface{}, ep *error) {
	d, e := i.GetData()
	if e == nil {
		*dp = d
	} else if e != nil && ep != nil && *ep == nil {
		*ep = e
	}

}

func (i *Wrapper) Error() (e error) {
	if i.err != nil {
		e = *i.err
	}
	return
}

func (i *Wrapper) GetDataOrNil() interface{} { return i.getData() }

func NewWrapper(i interface{}) (rd *Wrapper) {
	rd = new(Wrapper)
	rd.dat = i
	rd.access = make([]string, 1, 1)
	rd.access[0] = "<root>"
	return rd
}

func NewDictionary() *Wrapper {
	m := make(map[string]interface{})
	return NewWrapper(m)
}

func NewArray(l int) *Wrapper {
	m := make([]interface{}, l)
	return NewWrapper(m)
}

func NewNil() *Wrapper {
	return NewWrapper(nil)
}

func NewInt(i int) *Wrapper {
	return NewWrapper(i)
}

func NewInt64(i int64) *Wrapper {
	return NewWrapper(i)
}

func NewFloat64(f float64) *Wrapper {
	return NewWrapper(f)
}

func NewUint64(u uint64) *Wrapper {
	return NewWrapper(u)
}

func NewString(s string) *Wrapper {
	return NewWrapper(s)
}

func NewBool(b bool) *Wrapper {
	return NewWrapper(b)
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

func isFloat(v reflect.Value) bool {
	k := v.Kind()
	return k == reflect.Float32 || k == reflect.Float64
}

func (i *Wrapper) AccessPath() string {
	return strings.Join(i.access, "")
}

func (rd *Wrapper) GetFloat() (ret float64, err error) {
	if rd.err != nil {
		err = rd.err
	} else {
		v := reflect.ValueOf(rd.dat)
		if isFloat(v) {
			ret = float64(v.Float())
		} else if isInt(v) {
			ret = float64(v.Int())
		} else if isUint(v) {
			ret = float64(v.Uint())
		} else {
			err = rd.wrongType("float-like", v.Kind())
		}
	}
	return
}

func (w *Wrapper) GetFloatVoid(fp *float64, errp *error) {
	f, e := w.GetFloat()
	if e == nil {
		*fp = f
	} else if e != nil && errp != nil && *errp == nil {
		*errp = e
	}
}

func (rd *Wrapper) GetInt64() (ret int64, err error) {
	if rd.err != nil {
		err = rd.err
	} else {
		v := reflect.ValueOf(rd.dat)
		if isInt(v) {
			ret = v.Int()
		} else if isFloat(v) {
			ret = int64(v.Float())
		} else if !isUint(v) {
			err = rd.wrongType("int", v.Kind())
		} else if v.Uint() <= (1<<63 - 1) {
			ret = int64(v.Uint())
		} else {
			err = rd.NewError("Signed int64 overflow error")
		}
	}
	return
}

func (w *Wrapper) GetInt64Void(ip *int64, errp *error) {
	i, e := w.GetInt64()
	if e == nil {
		*ip = i
	} else if e != nil && errp != nil && *errp == nil {
		*errp = e
	}
}

func (rd *Wrapper) GetInt() (i int, err error) {
	i64, e := rd.GetInt64()
	return int(i64), e
}

func (w *Wrapper) GetIntVoid(ip *int, errp *error) {
	i, e := w.GetInt()
	if e == nil {
		*ip = i
	} else if e != nil && errp != nil && *errp == nil {
		*errp = e
	}
}

func (rd *Wrapper) GetUint() (u uint, err error) {
	u64, e := rd.GetUint64()
	return uint(u64), e
}

func (w *Wrapper) GetUintVoid(ip *uint, errp *error) {
	i, e := w.GetUint()
	if e == nil {
		*ip = i
	} else if e != nil && errp != nil && *errp == nil {
		*errp = e
	}
}

func (rd *Wrapper) GetUint64() (ret uint64, err error) {
	if rd.err != nil {
		err = rd.err
	} else {
		underflow := false
		v := reflect.ValueOf(rd.dat)
		if isUint(v) {
			ret = v.Uint()
		} else if isFloat(v) {
			if v.Float() <= 0 {
				underflow = true
			} else {
				ret = uint64(v.Float())
			}
		} else if !isInt(v) {
			err = rd.wrongType("uint", v.Kind())
		} else if v.Int() >= 0 {
			ret = uint64(v.Int())
		} else {
			underflow = true
		}

		if underflow {
			err = rd.NewError("Unsigned uint64 underflow error")

		}
	}
	return
}

func (w *Wrapper) GetUint64Void(ip *uint64, errp *error) {
	i, e := w.GetUint64()
	if e == nil {
		*ip = i
	} else if e != nil && errp != nil && *errp == nil {
		*errp = e
	}
}

func (rd *Wrapper) GetBool() (ret bool, err error) {
	if rd.err != nil {
		err = rd.err
	} else {
		v := reflect.ValueOf(rd.dat)
		k := v.Kind()
		if k == reflect.Bool {
			ret = v.Bool()
		} else {
			err = rd.wrongType("bool", k)
		}
	}
	return
}

func (w *Wrapper) GetBoolVoid(bp *bool, errp *error) {
	b, e := w.GetBool()
	if e == nil {
		*bp = b
	} else if e != nil && errp != nil && *errp == nil {
		*errp = e
	}
}

func (rd *Wrapper) GetString() (ret string, err error) {
	if rd.err != nil {
		err = rd.err
	} else if v := reflect.ValueOf(rd.dat); v.Kind() == reflect.String {
		ret = v.String()
	} else if b, ok := rd.dat.([]uint8); ok {
		ret = string(b)
	} else if b, ok := rd.dat.([]byte); ok {
		ret = string(b)
	} else {
		err = rd.wrongType("string", v.Kind())
	}
	return
}

func (w *Wrapper) GetStringVoid(sp *string, errp *error) {
	s, e := w.GetString()
	if e == nil {
		*sp = s
	} else if e != nil && errp != nil && *errp == nil {
		*errp = e
	}
}

func (rd *Wrapper) AtIndex(i int) *Wrapper {
	ret, v := rd.asArray()
	if v == nil {

	} else if len(v) <= i {
		ret.err = rd.NewError("index out of bounds %d >= %d", i, len(v))
	} else {
		ret.dat = v[i]
	}
	ret.access = append(ret.access, fmt.Sprintf("[%d]", i))
	return ret
}

func (rd *Wrapper) Len() (ret int, err error) {
	tmp, v := rd.asArray()
	if v == nil {
		err = tmp.err
	} else {
		ret = len(v)
	}
	return
}

func (i *Wrapper) Keys() (v []string, err error) {
	tmp, d := i.asDictionary()
	if d == nil {
		err = tmp.err
	} else {
		v = make([]string, len(d))
		var i int = 0
		for k, _ := range d {
			v[i] = k
			i++
		}
	}
	return
}

func (i *Wrapper) asArray() (ret *Wrapper, v []interface{}) {
	if i.err != nil {
		ret = i
	} else {
		var ok bool
		v, ok = (i.dat).([]interface{})
		ret = new(Wrapper)
		ret.access = i.access
		if !ok {
			ret.err = i.wrongType("array", reflect.ValueOf(i.dat).Kind())
		}
	}
	return
}

func (rd *Wrapper) IsNil() bool {
	return rd.dat == nil
}

func (rd *Wrapper) AtKey(s string) *Wrapper {
	ret, d := rd.asDictionary()

	if d != nil {
		val, found := d[s]
		if found {
			ret.dat = val
		} else {
			ret.dat = nil
		}
	}
	ret.access = append(ret.access, fmt.Sprintf(".%s", s))
	return ret
}

func (rd *Wrapper) ToDictionary() (out *Wrapper, e error) {
	tmp,_ := rd.asDictionary()
	if tmp.err != nil {
		e = tmp.err;
	} else {
		out = rd;
	}
	return
}

func (rd *Wrapper) ToArray() (out *Wrapper, e error) {
	tmp,_ := rd.asArray ()
	if tmp.err != nil {
		e = tmp.err;
	} else {
		out = rd;
	}
	return
}

func (w *Wrapper) SetKey(s string, val *Wrapper) error {
	b, d := w.asDictionary()
	if d != nil {
		d[s] = val.getData()
	}
	return b.Error()
}

func (w *Wrapper) SetIndex(i int, val *Wrapper) error {
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
		d, ok = (i.dat).(map[string]interface{})
		ret = new(Wrapper)
		ret.access = i.access
		if !ok {
			ret.err = i.wrongType("dict", reflect.ValueOf(i.dat).Kind())
		}
	}
	return
}

func (w *Wrapper) AtPath(path string) (ret *Wrapper) {
	bits := strings.Split(path, ".")
	ret = w
	for _, bit := range bits {
		if len(bit) > 0 && (bit[0] >= '0' && bit[0] <= '9') {
			// this is probably an int, use AtIndex instead
			if val, e := strconv.Atoi(bit); e == nil {
				ret = ret.AtIndex(val)
			} else {
				ret = ret.AtKey(bit)
			}
		} else if len(bit) > 0 {
			ret = ret.AtKey(bit)
		} else {
			break
		}

		if ret.dat == nil || ret.err != nil {
			break
		}
	}
	return ret
}
