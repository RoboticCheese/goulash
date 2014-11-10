package common

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func Test_Supermarketer_1(t *testing.T) {
	// Doesn't do anything just yet
}

func Test_Equals_1_Equal(t *testing.T) {
	c1 := Component{Endpoint: "somewhere"}
	c2 := Component{Endpoint: "somewhere"}
	res := Equals(&c1, &c2)
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_Equals_2_NotEqual(t *testing.T) {
	c1 := Component{Endpoint: "somewhere"}
	c2 := Component{Endpoint: "elsewhere"}
	res := Equals(&c1, &c2)
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Diff_1_Equal(t *testing.T) {
	c1 := Component{Endpoint: "abc", ETag: "def"}
	c2 := Component{Endpoint: "abc", ETag: "def"}
	pos1, neg1 := Diff(&c1, &c2, &Component{}, &Component{})
	pos2, neg2 := Diff(&c2, &c1, &Component{}, &Component{})
	for _, i := range []Supermarketer{pos1, neg1, pos2, neg2} {
		if i != nil {
			t.Fatalf("Expected nil, got: %v", i)
		}
	}
}

func Test_Diff_2_AddedAndDeletedData(t *testing.T) {
	c1 := Component{}
	c2 := Component{Endpoint: "abc", ETag: "def"}
	pos1, neg1 := Diff(&c1, &c2, &Component{}, &Component{})
	pos2, neg2 := Diff(&c2, &c1, &Component{}, &Component{})
	for _, i := range []Supermarketer{neg1, pos2} {
		if i != nil {
			t.Fatalf("Expected nil, got: %v", i)
		}
	}
	for _, k := range [][]Supermarketer{
		{pos1, &Component{Endpoint: "abc", ETag: "def"}},
		{neg2, &Component{Endpoint: "abc", ETag: "def"}},
	} {
		if !Equals(k[0], k[1]) {
			t.Fatalf("Expected %v, got: %v", k[1], k[0])
		}
	}
}

func Test_Diff_3_ChangedData(t *testing.T) {
	c1 := Component{Endpoint: "abc", ETag: "def"}
	c2 := Component{Endpoint: "uvw", ETag: "xyz"}
	pos1, neg1 := Diff(&c1, &c2, &Component{}, &Component{})
	pos2, neg2 := Diff(&c2, &c1, &Component{}, &Component{})
	for _, k := range [][]Supermarketer{
		{pos1, &Component{Endpoint: "uvw", ETag: "xyz"}},
		{neg1, &Component{Endpoint: "abc", ETag: "def"}},
		{pos2, &Component{Endpoint: "abc", ETag: "def"}},
		{neg2, &Component{Endpoint: "uvw", ETag: "xyz"}},
	} {
		if !Equals(k[0], k[1]) {
			t.Fatalf("Expected %v, got: %v", k[1], k[0])
		}
	}
}

func Test_emptyValue_1_EmptyString(t *testing.T) {
	res := emptyValue(reflect.ValueOf(""))
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_emptyValue_2_NonEmptyString(t *testing.T) {
	res := emptyValue(reflect.ValueOf("abc"))
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_emptyValue_3_EmptyPtr(t *testing.T) {
	c := Component{}
	res := emptyValue(reflect.ValueOf(&c))
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_emptyValue_4_NonEmptyPtr(t *testing.T) {
	c := Component{Endpoint: "abc"}
	res := emptyValue(reflect.ValueOf(&c))
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_emptyValue_5_EmptyMap(t *testing.T) {
	res := emptyValue(reflect.ValueOf(map[string]string{}))
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_emptyValue_6_NonEmptyMap(t *testing.T) {
	c := Component{Endpoint: "abc"}
	res := emptyValue(reflect.ValueOf(map[string]*Component{"thing": &c}))
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}
