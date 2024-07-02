package simplemon

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/status", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 500, rr.Code)

	var resp response
	assert.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	assert.Equal(t, bad, resp.Status)
}
