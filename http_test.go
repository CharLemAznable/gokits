package gokits

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFormIntValue(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			intValue1 := FormIntValue(r, "intValue1")
			if 1 != intValue1 {
				t.Errorf("intValue1 Should be 1")
			}
			intValue2 := FormIntValue(r, "intValue2")
			if 0 != intValue2 {
				t.Errorf("intValue2 Should be 0")
			}
			intValue3 := FormIntValueDefault(r, "intValue3", 4)
			if 3 != intValue3 {
				t.Errorf("intValue3 Should be 3")
			}
			intValue4 := FormIntValueDefault(r, "intValue4", 4)
			if 4 != intValue4 {
				t.Errorf("intValue4 Should be 4")
			}
			w.WriteHeader(http.StatusOK)
		}))
	_, err := NewHttpReq(testServer.URL).
		Params("intValue1", "1").
		Params("intValue2", "two").
		Params("intValue3", "3").Get()
	if nil != err {
		t.Errorf("Should has no error")
	}
}

func TestResponse(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ResponseJson(w, "{\"json\":\"JSON\"}")
		}))
	json, _ := NewHttpReq(testServer.URL).Get()
	if json != "{\"json\":\"JSON\"}" {
		t.Errorf("Should response {\"json\":\"JSON\"}")
	}

	testServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ResponseText(w, "plain text")
		}))
	text, _ := NewHttpReq(testServer.URL).Get()
	if text != "plain text" {
		t.Errorf("Should response plain text")
	}

	testServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ResponseHtml(w, "<html></html>")
		}))
	html, _ := NewHttpReq(testServer.URL).Get()
	if html != "<html></html>" {
		t.Errorf("Should response <html></html>")
	}
}

func TestResponseError(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ResponseErrorJson(w, http.StatusInternalServerError, "{\"json\":\"JSON\"}")
		}))
	json, _ := NewHttpReq(testServer.URL).Get()
	if json != "{\"json\":\"JSON\"}" {
		t.Errorf("Should response {\"json\":\"JSON\"}")
	}

	testServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ResponseErrorText(w, http.StatusInternalServerError, "plain text")
		}))
	text, _ := NewHttpReq(testServer.URL).Get()
	if text != "plain text" {
		t.Errorf("Should response plain text")
	}

	testServer = httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ResponseErrorHtml(w, http.StatusInternalServerError, "<html></html>")
		}))
	html, _ := NewHttpReq(testServer.URL).Get()
	if html != "<html></html>" {
		t.Errorf("Should response <html></html>")
	}
}
