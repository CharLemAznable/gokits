package gokits

import (
    "github.com/stretchr/testify/assert"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestFormIntValue(t *testing.T) {
    a := assert.New(t)
    testServer := httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            intValue1 := FormIntValue(r, "intValue1")
            a.Equal(1, intValue1)
            intValue2 := FormIntValue(r, "intValue2")
            a.Equal(0, intValue2)
            intValue3 := FormIntValueDefault(r, "intValue3", 4)
            a.Equal(3, intValue3)
            intValue4 := FormIntValueDefault(r, "intValue4", 4)
            a.Equal(4, intValue4)
            w.WriteHeader(http.StatusOK)
        }))
    _, err := NewHttpReq(testServer.URL).
        Params("intValue1", "1").
        Params("intValue2", "two").
        Params("intValue3", "3").Get()
    a.Nil(err)
}

func TestResponse(t *testing.T) {
    a := assert.New(t)
    testServer := httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            ResponseJson(w, "{\"json\":\"JSON\"}")
        }))
    json, _ := NewHttpReq(testServer.URL).Get()
    a.Equal("{\"json\":\"JSON\"}", json)

    testServer = httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            ResponseText(w, "plain text")
        }))
    text, _ := NewHttpReq(testServer.URL).Get()
    a.Equal("plain text", text)

    testServer = httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            ResponseHtml(w, "<html></html>")
        }))
    html, _ := NewHttpReq(testServer.URL).Get()
    a.Equal("<html></html>", html)
}

func TestResponseError(t *testing.T) {
    a := assert.New(t)
    testServer := httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            ResponseErrorJson(w, http.StatusInternalServerError, "{\"json\":\"JSON\"}")
        }))
    code, json, _ := NewHttpReq(testServer.URL).testGet()
    a.Equal(http.StatusInternalServerError, code)
    a.Equal("{\"json\":\"JSON\"}", json)

    testServer = httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            ResponseErrorText(w, http.StatusInternalServerError, "plain text")
        }))
    code, text, _ := NewHttpReq(testServer.URL).testGet()
    a.Equal(http.StatusInternalServerError, code)
    a.Equal("plain text", text)

    testServer = httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            ResponseErrorHtml(w, http.StatusInternalServerError, "<html></html>")
        }))
    code, html, _ := NewHttpReq(testServer.URL).testGet()
    a.Equal(http.StatusInternalServerError, code)
    a.Equal("<html></html>", html)
}
