package goulash

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
	g, err := New(ts.URL)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if g.Endpoint != ts.URL {
		t.Fatalf("Expected Endpoint: %v, got: %v", ts.URL, g.Endpoint)
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

func Test_New_2_404Error(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(http.NotFound))
	defer ts.Close()

	_, err := New(ts.URL)
	if err == nil {
		t.Fatalf("Expected an error but didn't get one")
	}
}
