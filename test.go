// -*- mode: go; tab-width: 4; c-basic-offset: 4; indent-tabs-mode: nil; -*-

package jsonw

import "testing"

func TestInt (t *testing.T) {
	const x = 100;
	w := NewInt(x);
	if v, _ := w.GetInt(); v != x {
		t.Errorf("%d != %d in GetInt() test", v, x);
	}
}
