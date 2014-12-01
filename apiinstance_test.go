package goulash

import (
	"net/http"
	"testing"
)

func TestNewAPIInstanceNoError(t *testing.T) {
	ts := StartHTTP("", nil)
	defer ts.Close()
	i, err := NewAPIInstance(ts.URL)
	for _, i := range [][]interface{}{
		{err, nil},
		{i.BaseURL, ts.URL},
		{i.Endpoint, ts.URL + "/api/v1"},
		{i.Version, "1"},
	} {
		if i[0] != i[1] {
			t.Fatalf("Expected: %v, got: %v", i[1], i[0])
		}
	}
}

func TestNewAPIInstanceConnError(t *testing.T) {
	ts := StartHTTP("", nil)
	ts.Close()

	_, err := NewAPIInstance(ts.URL)
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func TestNewAPIInstance404Error(t *testing.T) {
	ts := StartHTTP(http.NotFound, nil)
	defer ts.Close()

	_, err := NewAPIInstance(ts.URL)
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func TestNewAPIInstanceRealData(t *testing.T) {
	_, err := NewAPIInstance("https://supermarket.getchef.com")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func TestInitAPIInstanceEmptyResult(t *testing.T) {
	i := InitAPIInstance()
	for _, k := range []string{
		i.Endpoint,
		i.ETag,
		i.BaseURL,
		i.Version,
	} {
		if k != "" {
			t.Fatalf("Expected empty string, got: %v", k)
		}
	}
}

func TestAPIInstanceEmptyEmpty(t *testing.T) {
	a := new(APIInstance)
	res := a.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func TestAPIInstanceEmptyHasEndpoint(t *testing.T) {
	a := new(APIInstance)
	a.Endpoint = "https://example.com"
	res := a.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestAPIInstanceEmptyHasBaseURL(t *testing.T) {
	a := new(APIInstance)
	a.BaseURL = "https://example.com"
	res := a.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func TestAPIInstanceEmptyHasVersion(t *testing.T) {
	a := new(APIInstance)
	a.Version = "1"
	res := a.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}
