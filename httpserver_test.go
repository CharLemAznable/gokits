package gokits

import (
    "github.com/stretchr/testify/assert"
    "net/http"
    "net/http/httptest"
    "net/url"
    "testing"
)

func TestReverseProxy(t *testing.T) {
    a := assert.New(t)
    testServer := httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
        }))
    testServerURL, _ := url.Parse(testServer.URL)
    reverseProxy := ReverseProxy(testServerURL)
    reverseServer := httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            reverseProxy.ServeHTTP(w, r)
        }))
    code, _, _ := NewHttpReq(reverseServer.URL).testGet()
    a.Equal(http.StatusOK, code)
}

func TestHandleFuncOptions(t *testing.T) {
    a := assert.New(t)
    options := defaultHandleFuncOptions
    for _, o := range []HandleFuncOption{DumpRequestDisabled,
        GzipResponseDisabled, ModelContextDisabled, ContextPathDisabled} {
        o(&options)
    }
    a.False(options.DumpRequestEnabled || options.GzipResponseEnabled ||
        options.ModelContextEnabled || options.ContextPathEnabled)

    options = defaultHandleFuncOptions
    for _, o := range []HandleFuncOption{DumpRequestEnabled,
        GzipResponseEnabled, ModelContextEnabled, ContextPathEnabled} {
        o(&options)
    }
    a.False(!options.DumpRequestEnabled || !options.GzipResponseEnabled ||
        !options.ModelContextEnabled || !options.ContextPathEnabled)
}

func TestHandleFunc(t *testing.T) {
    a := assert.New(t)
    mux := http.NewServeMux()
    HandleFunc(mux, "/index",
        func(w http.ResponseWriter, r *http.Request) {
            ResponseText(w, "index")
        })
    testServer := httptest.NewServer(mux)
    code, resp, _ := NewHttpReq(testServer.URL + "/index").testGet()
    a.Equal(http.StatusOK, code)
    a.Equal("index", resp)
}
