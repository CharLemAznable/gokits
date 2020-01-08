package gokits

import (
    "net/http"
    "net/http/httptest"
    "net/url"
    "testing"
)

func TestReverseProxy(t *testing.T) {
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
    if code != http.StatusOK {
        t.Errorf("Should response http.StatusOK")
    }
}

func TestHandleFuncOptions(t *testing.T) {
    options := defaultHandleFuncOptions
    for _, o := range []HandleFuncOption{DumpRequestDisabled,
        GzipResponseDisabled, ModelContextDisabled, ContextPathDisabled} {
        o(&options)
    }
    if options.DumpRequestEnabled || options.GzipResponseEnabled ||
        options.ModelContextEnabled || options.ContextPathEnabled {
        t.Errorf("Should disabled")
    }

    options = defaultHandleFuncOptions
    for _, o := range []HandleFuncOption{DumpRequestEnabled,
        GzipResponseEnabled, ModelContextEnabled, ContextPathEnabled} {
        o(&options)
    }
    if !options.DumpRequestEnabled || !options.GzipResponseEnabled ||
        !options.ModelContextEnabled || !options.ContextPathEnabled {
        t.Errorf("Should enabled")
    }
}

func TestHandleFunc(t *testing.T) {
    mux := http.NewServeMux()
    HandleFunc(mux, "/index",
        func(w http.ResponseWriter, r *http.Request) {
            ResponseText(w, "index")
        })
    testServer := httptest.NewServer(mux)
    code, resp, _ := NewHttpReq(testServer.URL + "/index").testGet()
    if code != http.StatusOK {
        t.Errorf("Should response http.StatusOK")
    }
    if resp != "index" {
        t.Errorf("Should response index")
    }
}
