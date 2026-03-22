package respond_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"fintrack/pkg/respond"
)

func TestJSON_SetsContentType(t *testing.T) {
	rr := httptest.NewRecorder()
	respond.JSON(rr, http.StatusOK, map[string]string{"hello": "world"})

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", ct)
	}
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestSuccess_WrapsData(t *testing.T) {
	rr := httptest.NewRecorder()
	respond.Success(rr, http.StatusCreated, map[string]string{"id": "123"})

	var body map[string]any
	json.NewDecoder(rr.Body).Decode(&body)

	if body["success"] != true {
		t.Error("expected success=true")
	}
	data, ok := body["data"].(map[string]any)
	if !ok {
		t.Fatal("expected data field to be an object")
	}
	if data["id"] != "123" {
		t.Errorf("expected id=123, got %v", data["id"])
	}
}

func TestError_WrapsMessage(t *testing.T) {
	rr := httptest.NewRecorder()
	respond.Error(rr, http.StatusNotFound, "resource not found")

	var body map[string]any
	json.NewDecoder(rr.Body).Decode(&body)

	if body["success"] != false {
		t.Error("expected success=false")
	}
	if body["error"] != "resource not found" {
		t.Errorf("unexpected error message: %v", body["error"])
	}
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestValidationError_IncludesFields(t *testing.T) {
	rr := httptest.NewRecorder()
	respond.ValidationError(rr, map[string]string{
		"email":    "required",
		"password": "must be at least 8 characters",
	})

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", rr.Code)
	}

	var body map[string]any
	json.NewDecoder(rr.Body).Decode(&body)

	fields, ok := body["fields"].(map[string]any)
	if !ok {
		t.Fatal("expected fields in response")
	}
	if fields["email"] != "required" {
		t.Errorf("expected email field error, got %v", fields["email"])
	}
}
