package apiinstance

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func start_http() (ts *httptest.Server) {
	response := "Everything's okay!"
	ts = httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, response)
			},
		),
	)
	return
}

func Test_New_1_NoError(t *testing.T) {
	ts := start_http()
	defer ts.Close()
	i, err := New(ts.URL)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if i.BaseURL != ts.URL {
		t.Fatalf("Expected BaseURL: %v, got: %v", ts.URL, i.BaseURL)
	}
	endpoint := ts.URL + "/api/v1"
	if i.Endpoint != endpoint {
		t.Fatalf("Expected Endpoint: %v, got: %v", endpoint, i.Endpoint)
	}
	if i.Version != "1" {
		t.Fatalf("Expected Version: 1, got: %v", i.Version)
	}
}

func Test_New_2_ConnError(t *testing.T) {
	ts := start_http()
	ts.Close()

	_, err := New(ts.URL)
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func Test_New_3_404Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(http.NotFound))
	defer ts.Close()

	_, err := New(ts.URL)
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}

func Test_New_4_RealData(t *testing.T) {
	_, err := New("https://supermarket.getchef.com")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
}

func Test_Empty_1_Empty(t *testing.T) {
	a := new(APIInstance)
	res := a.Empty()
	if res != true {
		t.Fatalf("Expected true, got: %v", res)
	}
}

func Test_Empty_2_HasEndpoint(t *testing.T) {
	a := new(APIInstance)
	a.Endpoint = "https://example.com"
	res := a.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Empty_3_HasBaseURL(t *testing.T) {
	a := new(APIInstance)
	a.BaseURL = "https://example.com"
	res := a.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}

func Test_Empty_4_HasVersion(t *testing.T) {
	a := new(APIInstance)
	a.Version = "1"
	res := a.Empty()
	if res != false {
		t.Fatalf("Expected false, got: %v", res)
	}
}
