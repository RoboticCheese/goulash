package common

import (
	"reflect"
	"testing"
)

type thing struct {
	Endpoint string
	ETag     string
}

func (t *thing) Empty() (empty bool) {
	if t.Endpoint != "" || t.ETag != "" {
		empty = false
	} else {
		empty = true
	}
	return
}

func TestEqualsEqual(t *testing.T) {
	c1 := thing{Endpoint: "somewhere"}
	c2 := thing{Endpoint: "somewhere"}
	res := Equals(&c1, &c2)
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func TestEqualsNotEqual(t *testing.T) {
	c1 := thing{Endpoint: "somewhere"}
	c2 := thing{Endpoint: "elsewhere"}
	res := Equals(&c1, &c2)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestDiffEqual(t *testing.T) {
	c1 := thing{Endpoint: "abc", ETag: "def"}
	c2 := thing{Endpoint: "abc", ETag: "def"}
	pos1, neg1 := Diff(&c1, &c2, &thing{}, &thing{})
	pos2, neg2 := Diff(&c2, &c1, &thing{}, &thing{})
	for _, i := range []Supermarketer{pos1, neg1, pos2, neg2} {
		if i != nil {
			t.Fatalf("Expected nil, got: %v", i)
		}
	}
}

func TestDiffAddedAndDeletedData(t *testing.T) {
	c1 := thing{}
	c2 := thing{Endpoint: "abc", ETag: "def"}
	pos1, neg1 := Diff(&c1, &c2, &thing{}, &thing{})
	pos2, neg2 := Diff(&c2, &c1, &thing{}, &thing{})
	for _, i := range []Supermarketer{neg1, pos2} {
		if i != nil {
			t.Fatalf("Expected nil, got: %v", i)
		}
	}
	for _, k := range [][]Supermarketer{
		{pos1, &thing{Endpoint: "abc", ETag: "def"}},
		{neg2, &thing{Endpoint: "abc", ETag: "def"}},
	} {
		if !Equals(k[0], k[1]) {
			t.Fatalf("Expected %v, got: %v", k[1], k[0])
		}
	}
}

func TestDiffChangedData(t *testing.T) {
	c1 := thing{Endpoint: "abc", ETag: "def"}
	c2 := thing{Endpoint: "uvw", ETag: "xyz"}
	pos1, neg1 := Diff(&c1, &c2, &thing{}, &thing{})
	pos2, neg2 := Diff(&c2, &c1, &thing{}, &thing{})
	for _, k := range [][]Supermarketer{
		{pos1, &thing{Endpoint: "uvw", ETag: "xyz"}},
		{neg1, &thing{Endpoint: "abc", ETag: "def"}},
		{pos2, &thing{Endpoint: "abc", ETag: "def"}},
		{neg2, &thing{Endpoint: "uvw", ETag: "xyz"}},
	} {
		if !Equals(k[0], k[1]) {
			t.Fatalf("Expected %v, got: %v", k[1], k[0])
		}
	}
}

func TestemptyValueEmptyString(t *testing.T) {
	res := emptyValue(reflect.ValueOf(""))
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func TestemptyValueNonEmptyString(t *testing.T) {
	res := emptyValue(reflect.ValueOf("abc"))
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestemptyValueEmptyPtr(t *testing.T) {
	c := thing{}
	res := emptyValue(reflect.ValueOf(&c))
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func TestemptyValueNonEmptyPtr(t *testing.T) {
	c := thing{Endpoint: "abc"}
	res := emptyValue(reflect.ValueOf(&c))
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestemptyValueEmptyMap(t *testing.T) {
	res := emptyValue(reflect.ValueOf(map[string]string{}))
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func TestemptyValueNonEmptyMap(t *testing.T) {
	c := thing{Endpoint: "abc"}
	res := emptyValue(reflect.ValueOf(map[string]*thing{"thing": &c}))
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}
