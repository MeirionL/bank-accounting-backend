package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondWithJSON(t *testing.T) {
	testResponse := httptest.NewRecorder()
	payload := map[string]interface{}{
		"message": "foo",
	}

	respondWithJSON(testResponse, http.StatusOK, payload)

	if testResponse.Code != http.StatusOK {
		t.Errorf("expected status %v, but got %v", http.StatusOK, testResponse.Code)
	}

	contentType := testResponse.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected content type header application/json but got %s", contentType)
	}

	var responsePayload map[string]interface{}
	err := json.Unmarshal(testResponse.Body.Bytes(), &responsePayload)
	if err != nil {
		t.Errorf("failed to unmarshal JSON response: %v", err)
	}

	message, ok := responsePayload["message"].(string)
	if !ok {
		t.Error("expected message key in JSON payload but it's missing or not a string")
	}

	if message != "foo" {
		t.Errorf("expected message foo but got %s", message)
	}
}
