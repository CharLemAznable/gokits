package gokits

import (
    "net/http"
    "net/http/httputil"
    "net/url"
)

var ReverseProxyUserAgent = ""

func ReverseProxy(target *url.URL) *httputil.ReverseProxy {
    targetQuery := target.RawQuery
    director := func(req *http.Request) {
        req.Host = target.Host // Different from the default NewSingleHostReverseProxy()

        req.URL.Scheme = target.Scheme
        req.URL.Host = target.Host
        req.URL.Path = PathJoin(target.Path, req.URL.Path)
        if targetQuery == "" || req.URL.RawQuery == "" {
            req.URL.RawQuery = targetQuery + req.URL.RawQuery
        } else {
            req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
        }
        if _, ok := req.Header["User-Agent"]; !ok {
            req.Header.Set("User-Agent", ReverseProxyUserAgent)
        }
    }
    return &httputil.ReverseProxy{Director: director}
}

type HttpServerConfig struct {
    Port        int
    ContextPath string
}

var GlobalHttpServerConfig *HttpServerConfig = nil

type HandleFuncOptions struct {
    DumpRequestEnabled  bool
    GzipResponseEnabled bool
    ModelContextEnabled bool
    ContextPathEnabled  bool
}

var defaultHandleFuncOptions = HandleFuncOptions{
    DumpRequestEnabled:  true,
    GzipResponseEnabled: true,
    ModelContextEnabled: true,
    ContextPathEnabled:  true,
}

type HandleFuncOption func(*HandleFuncOptions)

var DumpRequestEnabled = func(o *HandleFuncOptions) { o.DumpRequestEnabled = true }
var DumpRequestDisabled = func(o *HandleFuncOptions) { o.DumpRequestEnabled = false }
var GzipResponseEnabled = func(o *HandleFuncOptions) { o.GzipResponseEnabled = true }
var GzipResponseDisabled = func(o *HandleFuncOptions) { o.GzipResponseEnabled = false }
var ModelContextEnabled = func(o *HandleFuncOptions) { o.ModelContextEnabled = true }
var ModelContextDisabled = func(o *HandleFuncOptions) { o.ModelContextEnabled = false }
var ContextPathEnabled = func(o *HandleFuncOptions) { o.ContextPathEnabled = true }
var ContextPathDisabled = func(o *HandleFuncOptions) { o.ContextPathEnabled = false }

func HandleFunc(mux *http.ServeMux, path string, handler http.HandlerFunc, opts ...HandleFuncOption) {
    options := defaultHandleFuncOptions
    for _, o := range opts {
        o(&options)
    }

    wrapHandler := handler
    if options.DumpRequestEnabled {
        wrapHandler = DumpRequest(wrapHandler)
    }
    if options.GzipResponseEnabled {
        wrapHandler = GzipResponse(wrapHandler)
    }
    if options.ModelContextEnabled {
        wrapHandler = ServeModelContext(wrapHandler)
    }

    wrapPath := path
    if options.ContextPathEnabled && nil != GlobalHttpServerConfig {
        wrapPath = PathJoin(GlobalHttpServerConfig.ContextPath, wrapPath)
    }

    mux.HandleFunc(wrapPath, wrapHandler)
}
