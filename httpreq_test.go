package gokits

import (
    "github.com/stretchr/testify/assert"
    "io/ioutil"
    "log"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHttpReq_Get(t *testing.T) {
    a := assert.New(t)
    testServer := httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            value := r.FormValue("key")
            a.Equal("value", value)
            w.WriteHeader(http.StatusOK)
        }))
    _, err := NewHttpReq(testServer.URL).Params("key", "value").Get()
    a.Nil(err)
}

func TestHttpReq_Post(t *testing.T) {
    a := assert.New(t)
    testServer := httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            value := r.FormValue("key")
            a.Equal("value", value)
            w.WriteHeader(http.StatusOK)
        }))
    _, err := NewHttpReq(testServer.URL).Params("key", "value").Post()
    a.Nil(err)
}

func TestHttpReq_Post_Body(t *testing.T) {
    a := assert.New(t)
    testServer := httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            body, _ := RequestBody(r)
            a.Equal("{\"key\":\"value\"}", body)
            a.Equal("application/json", r.Header.Get("Content-Type"))
            w.WriteHeader(http.StatusOK)
        }))
    _, err := NewHttpReq(testServer.URL).
        RequestBody("{\"key\":\"value\"}").
        Prop("Content-Type", "application/json").Post()
    a.Nil(err)
}

func (httpReq *HttpReq) testGet() (int, string, error) {
    request, err := httpReq.createGetRequest()
    if nil != err {
        return http.StatusBadRequest, "", err
    }
    httpReq.commonSettings(request)
    return httpReq.testDoRequest(request)
}

func (httpReq *HttpReq) testPost() (int, string, error) {
    request, err := httpReq.createPostRequest()
    if nil != err {
        return http.StatusBadRequest, "", err
    }
    httpReq.commonSettings(request)
    httpReq.postSettings(request)
    return httpReq.testDoRequest(request)
}

func (httpReq *HttpReq) testDoRequest(request *http.Request) (int, string, error) {
    httpReq.setHeaders(request)
    httpReq.addCookies(request)

    request.Close = true // fd leak without setting this
    response, err := http.DefaultClient.Do(request)
    if nil != err {
        log.Printf("%s: %s, STATUS CODE = %d\n\n%s\n",
            request.Method, request.URL.String(), Condition(nil != response,
                func() interface{} { return response.StatusCode }, -1), err.Error())
        return response.StatusCode, "", err
    }

    defer func() { _ = response.Body.Close() }()
    body, err := ioutil.ReadAll(response.Body)
    if nil != err {
        return response.StatusCode, "", err
    }
    return response.StatusCode, string(body), nil
}
