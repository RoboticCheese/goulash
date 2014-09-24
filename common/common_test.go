package common

import (
	"testing"
)

func Test_Supermarketer_1(t *testing.T) {
	// Doesn't do anything just yet
}

func Test_Component_1(t *testing.T) {
	type Thing struct {
		Component
	}

	res := new(Thing)
	res.Endpoint = "https://supermarket.getchef.com"
	if res.Endpoint != "https://supermarket.getchef.com" {
		t.Fatalf("Expected Endpoint to be set, got: %v", res.Endpoint)
	}
}
