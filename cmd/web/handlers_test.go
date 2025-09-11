package main

import (
	"bytes"
	"example/internal/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {

	// Initialize a new httpt.ResponseRecorder
	rr := httptest.NewRecorder()

	// Initialize a new dummy http.Request
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Call the ping handler function
	ping(rr, r)

	rs := rr.Result()

	// Check the status code written by the ping handler was 200
	assert.Equal(t, rs.StatusCode, http.StatusOK)

	// And we can check that the response body written by the ping handler aquals ok
	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
