package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestInfoRequestHandler(t *testing.T) {
	// switch to test mode to reduce 'noisy' output
	gin.SetMode(gin.TestMode)

	// setup the router, just like done in main function
	// and register the routes
	r := gin.Default()
	r.GET("/info", infoRequest_Handler)

	// create a fake request to test targeted function. Make sure the second argument
	// second argument must match the defined route in the router setup (above)
	test_request, err := http.NewRequest(http.MethodGet, "/info", nil)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}

	// create a response recorder to inspect the response
	w := httptest.NewRecorder()

	// perform the request using the recorder
	r.ServeHTTP(w, test_request)

	// check to see if the response was what you expected
	// first check status code
	if w.Code != http.StatusOK {
		t.Fatalf("Expected to get status %d but instead got %d\n", http.StatusOK, w.Code)
	}

	// next check content type
	if ctype := w.Header().Get("Content-Type"); ctype != "text/plain; charset=utf-8" {
		t.Fatalf("Expected to get code %s but instead got %s\n", "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
	}

}
