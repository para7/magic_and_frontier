package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	homeHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "連絡先フォーム") {
		t.Fatalf("body should include page title")
	}
}

func TestSubmitHandlerValidationError(t *testing.T) {
	form := url.Values{}
	form.Set("name", "")
	form.Set("phone", "abc")
	form.Set("mode", "birthdate")
	form.Set("birthdate", "")

	req := httptest.NewRequest(http.MethodPost, "/contact/submit", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	rr := httptest.NewRecorder()

	submitHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "名前は1文字以上で入力してください") {
		t.Fatalf("expected name validation error")
	}
	if !strings.Contains(body, "電話番号は数字のみで入力してください") {
		t.Fatalf("expected phone validation error")
	}
	if !strings.Contains(body, "生年月日を入力してください") {
		t.Fatalf("expected birthdate validation error")
	}
}

func TestSubmitHandlerSuccess(t *testing.T) {
	form := url.Values{}
	form.Set("name", "太郎")
	form.Set("phone", "09012345678")
	form.Set("mode", "birthdate")
	form.Set("birthdate", "2000-01-01")

	req := httptest.NewRequest(http.MethodPost, "/contact/submit", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("HX-Request", "true")
	rr := httptest.NewRecorder()

	submitHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	if !strings.Contains(rr.Body.String(), "送信しました。") {
		t.Fatalf("expected success message")
	}
}

func TestModeFieldsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/contact/mode-fields?mode=heightweight", nil)
	req.Header.Set("HX-Request", "true")
	rr := httptest.NewRecorder()

	modeFieldsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "id=\"height\"") {
		t.Fatalf("expected height field")
	}
	if !strings.Contains(body, "id=\"weight\"") {
		t.Fatalf("expected weight field")
	}
	if strings.Contains(body, "id=\"birthdate\"") {
		t.Fatalf("did not expect birthdate field")
	}
}
