package syncutils

import "testing"

func TestAtomicPtr(t *testing.T) {
	a := "a"
	b := "b"
	c := "c"

	var x AtomicPtr[string]

	if x.Load() != nil {
		t.Fatalf("x.Load() != nil")
	}

	x.ptr = &a
	if x.Load() != &a {
		t.Fatalf("x.Load() != &a")
	}

	x.Store(&b)
	if x.Load() != &b {
		t.Fatalf("x.Load() != &b")
	}

	x.Store(&c)
	if x.ptr != &c {
		t.Fatalf("x.ptr != &c")
	}
}
