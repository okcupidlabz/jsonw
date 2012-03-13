// -*- mode: go; tab-width: 4; c-basic-offset: 4; indent-tabs-mode: nil; -*-

package jsonw

import "testing"

func TestInt (t *testing.T) {
	const x = 100;
	w := NewInt(x);
	if v, _ := w.GetInt(); v != x  {
		t.Errorf("%d != %d in GetInt() test", v, x);
	}
}


func TestBigInt (t *testing.T) {
	const x = 1 << 62 + 55555;
	w := NewInt64(x);
	if v, _ := w.GetInt(); v != x  {
		t.Errorf("Big int test failed");
	}
    if v, _ := w.GetUint(); v != x {
		t.Errorf("Big uint test failed");
    }
}

func TestDict (t *testing.T) {
    w := NewDictionary()
    const dog = 3333
    var cat string = "meow"
    
    w.SetKey("dog", NewInt(dog));
    w.SetKey("cat", NewString(cat));

    if v, _ := w.AtKey("dog").GetInt(); v != dog {
        t.Errorf("Dictionary fail for 'dog': %d != %d", v, dog);
    }

    if v, _ := w.AtKey("cat").GetString(); v != cat {
        t.Errorf("Dictionary fail for 'dog': %s != %s", v, cat);
    }

}

