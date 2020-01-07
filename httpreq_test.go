package gokits

import (
    "io/ioutil"
    "log"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHttpReq_Get(t *testing.T) {
    testServer := httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            value := r.FormValue("key")
            if "value" != value {
                t.Errorf("r.FormValue(\"key\") should be \"value\"")
            }
            w.WriteHeader(http.StatusOK)
        }))
    _, err := NewHttpReq(testServer.URL).Params("key", "value").Get()
    if nil != err {
        t.Errorf("Should has no error")
    }
}

func TestHttpReq_Post(t *testing.T) {
    testServer := httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            value := r.FormValue("key")
            if "value" != value {
                t.Errorf("r.FormValue(\"key\") should be \"value\"")
            }
            w.WriteHeader(http.StatusOK)
        }))
    _, err := NewHttpReq(testServer.URL).Params("key", "value").Post()
    if nil != err {
        t.Errorf("Should has no error")
    }
}

func TestHttpReq_Post_Body(t *testing.T) {
    testServer := httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            body, _ := RequestBody(r)
            if "{\"key\":\"value\"}" != body {
                t.Errorf("requestBody should be {\"key\":\"value\"}")
            }
            if "application/json" != r.Header.Get("Content-Type") {
                t.Errorf("contentType should be application/json")
            }
            w.WriteHeader(http.StatusOK)
        }))
    _, err := NewHttpReq(testServer.URL).
        RequestBody("{\"key\":\"value\"}").
        Prop("Content-Type", "application/json").Post()
    if nil != err {
        t.Errorf("Should has no error")
    }
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
