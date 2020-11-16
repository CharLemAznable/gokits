package gokits

import (
    "compress/gzip"
    "context"
    "fmt"
    "github.com/kataras/golog"
    "github.com/tdewolff/minify"
    "github.com/tdewolff/minify/css"
    "github.com/tdewolff/minify/html"
    "github.com/tdewolff/minify/js"
    "io"
    "net/http"
    "net/http/httputil"
    "strings"
)

func DumpRequest(handlerFunc http.HandlerFunc) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        // Save a copy of this request for debugging.
        requestDump, err := httputil.DumpRequest(request, true)
        if err != nil {
            golog.Error(err)
        }
        golog.Debug(string(requestDump))
        handlerFunc(writer, request)
    }
}

type GzipResponseWriter struct {
    io.Writer
    http.ResponseWriter
}

func (w GzipResponseWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
}

func GzipResponse(handlerFunc http.HandlerFunc) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        if !strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
            handlerFunc(writer, request)
            return
        }
        writer.Header().Set("Content-Encoding", "gzip")
        gz := gzip.NewWriter(writer)
        defer func() { _ = gz.Close() }()
        gzr := GzipResponseWriter{Writer: gz, ResponseWriter: writer}
        handlerFunc(gzr, request)
    }
}

type ModelCtx struct {
    context.Context
    Model map[string]interface{}
}

func (m *ModelCtx) String() string {
    return fmt.Sprintf("%v.WithModel(%#v)", m.Context, Json(m.Model))
}

func (m *ModelCtx) Value(key interface{}) interface{} {
    keyStr, ok := key.(string)
    if !ok {
        return m.Context.Value(key)
    }
    value, ok := m.Model[keyStr]
    if !ok {
        return m.Context.Value(key)
    }
    return value
}

func ModelContext(parent context.Context) *ModelCtx {
    switch parent.(type) {
    case *ModelCtx:
        return parent.(*ModelCtx)
    default:
        return &ModelCtx{parent, map[string]interface{}{}}
    }
}

func ModelContextWithValue(parent context.Context, key string, val interface{}) *ModelCtx {
    if "" == key {
        panic("empty key")
    }
    modelCtx := ModelContext(parent)
    modelCtx.Model[key] = val
    return modelCtx
}

func ServeModelContext(handlerFunc http.HandlerFunc) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        handlerFunc(writer, request.WithContext(
            ModelContext(request.Context())))
    }
}

type ModelCtxValueFunc func() (string, interface{})

func ServeModelContextWithValueFunc(handlerFunc http.HandlerFunc, valueFunc ModelCtxValueFunc) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        key, val := valueFunc()
        handlerFunc(writer, request.WithContext(
            ModelContextWithValue(request.Context(), key, val)))
    }
}

//noinspection GoUnusedExportedFunction
func ServeRedirect(redirect string) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        http.Redirect(writer, request, redirect, http.StatusFound)
    }
}

func ServeGet(handlerFunc http.HandlerFunc) http.HandlerFunc {
    return ServeMethod(handlerFunc, http.MethodGet)
}

func ServePost(handlerFunc http.HandlerFunc) http.HandlerFunc {
    return ServeMethod(handlerFunc, http.MethodPost)
}

func ServeMethod(handlerFunc http.HandlerFunc, httpMethods ...string) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        for _, method := range httpMethods {
            if method == request.Method {
                handlerFunc(writer, request)
                return
            }
        }
        writer.WriteHeader(http.StatusNotFound)
    }
}

func ServeAjax(handlerFunc http.HandlerFunc) http.HandlerFunc {
    return func(writer http.ResponseWriter, request *http.Request) {
        if !IsAjaxRequest(request) {
            writer.WriteHeader(http.StatusNotFound)
            return
        }
        handlerFunc(writer, request)
    }
}

func IsAjaxRequest(request *http.Request) bool {
    return "XMLHttpRequest" == request.Header.Get("X-Requested-With")
}

func MinifyHTML(htmlString string, devMode bool) string {
    if devMode {
        return htmlString
    }

    mini := minify.New()
    mini.AddFunc("text/html", html.Minify)
    minified, _ := mini.String("text/html", htmlString)
    return minified
}

func MinifyCSS(cssString string, devMode bool) string {
    if devMode {
        return cssString
    }

    mini := minify.New()
    mini.AddFunc("text/css", css.Minify)

    minifiedCSS, err := mini.String("text/css", cssString)
    if err != nil {
        fmt.Println("mini css:", err.Error())
    }

    return minifiedCSS
}

func MinifyJs(jsString string, devMode bool) string {
    if devMode {
        return jsString
    }

    mini := minify.New()
    mini.AddFunc("text/javascript", js.Minify)
    minifiedJs, err := mini.String("text/javascript", jsString)
    if err != nil {
        fmt.Println("mini js:", err.Error())
    }

    return minifiedJs
}

var EmptyHandler http.HandlerFunc = func(_ http.ResponseWriter, _ *http.Request) {}
