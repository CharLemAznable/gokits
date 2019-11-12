package gokits

import (
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "strings"
)

type HttpReq struct {
    baseUrl string
    req     string
    params  map[string]string
    body    string
    props   []prop
    cookies []*http.Cookie
}

type prop struct {
    name  string
    value string
}

func NewHttpReq(baseUrl string) *HttpReq {
    httpReq := new(HttpReq)
    httpReq.baseUrl = baseUrl
    httpReq.params = make(map[string]string)
    httpReq.props = make([]prop, 0)
    httpReq.cookies = make([]*http.Cookie, 0)
    return httpReq
}

func (httpReq *HttpReq) Req(req string) *HttpReq {
    httpReq.req = req
    return httpReq
}

func (httpReq *HttpReq) Params(name string, value string, more ...string) *HttpReq {
    if 0 != len(name) || 0 != len(value) {
        httpReq.params[name] = value
    }

    for i := 0; i < len(more); i += 2 {
        if i+1 >= len(more) {
            break
        }

        k, v := more[i], more[i+1]
        if 0 != len(k) || 0 != len(v) {
            httpReq.params[k] = v
        }
    }

    return httpReq
}

func (httpReq *HttpReq) ParamsMapping(params map[string]string) *HttpReq {
    if nil == params {
        return httpReq
    }

    for key, value := range params {
        if 0 != len(key) || 0 != len(value) {
            httpReq.params[key] = value
        }
    }
    return httpReq
}

func (httpReq *HttpReq) RequestBody(requestBody string) *HttpReq {
    httpReq.body = requestBody
    return httpReq
}

func (httpReq *HttpReq) Prop(name string, value string) *HttpReq {
    if 0 == len(name) || 0 == len(value) {
        return httpReq
    }
    httpReq.props = append(httpReq.props, prop{name: name, value: value})
    return httpReq
}

func (httpReq *HttpReq) Cookie(cookie *http.Cookie) *HttpReq {
    if 0 == len(cookie.Name) || 0 == len(cookie.Value) {
        return httpReq
    }
    httpReq.cookies = append(httpReq.cookies, cookie)
    return httpReq
}

func (httpReq *HttpReq) Get() (string, error) {
    request, err := httpReq.createGetRequest()
    if nil != err {
        return "", err
    }
    httpReq.commonSettings(request)
    return httpReq.doRequest(request)
}

func (httpReq *HttpReq) Post() (string, error) {
    request, err := httpReq.createPostRequest()
    if nil != err {
        return "", err
    }
    httpReq.commonSettings(request)
    httpReq.postSettings(request)
    return httpReq.doRequest(request)
}

func (httpReq *HttpReq) doRequest(request *http.Request) (string, error) {
    httpReq.setHeaders(request)
    httpReq.addCookies(request)

    request.Close = true // fd leak without setting this
    response, err := http.DefaultClient.Do(request)
    if nil != err {
        log.Printf("%s: %s, STATUS CODE = %d\n\n%s\n",
            request.Method, request.URL.String(), Condition(nil != response,
                func() interface{} { return response.StatusCode }, -1), err.Error())
        return "", err
    }

    defer func() { _ = response.Body.Close() }()
    body, err := ioutil.ReadAll(response.Body)
    if nil != err {
        return "", err
    }
    return string(body), nil
}

func (httpReq *HttpReq) createGetRequest() (*http.Request, error) {
    values := url.Values{}
    for key, value := range httpReq.params {
        values.Add(key, value)
    }
    encoded := values.Encode()

    urlStr := httpReq.baseUrl + httpReq.req
    if len(encoded) > 0 {
        urlStr = urlStr + "?" + encoded
    }
    return http.NewRequest("GET", urlStr, nil)
}

func (httpReq *HttpReq) createPostRequest() (*http.Request, error) {
    values := url.Values{}
    for key, value := range httpReq.params {
        values.Add(key, value)
    }
    encoded := values.Encode()

    if len(httpReq.body) > 0 {
        if len(encoded) > 0 {
            encoded = encoded + "&"
        }
        encoded = encoded + httpReq.body
    }

    urlStr := httpReq.baseUrl + httpReq.req
    bodyReader := strings.NewReader(encoded)
    return http.NewRequest("POST", urlStr, bodyReader)
}

func (httpReq *HttpReq) commonSettings(request *http.Request) {
    request.Header.Set("Accept-Charset", "UTF-8")
}

func (httpReq *HttpReq) setHeaders(request *http.Request) {
    for _, prop := range httpReq.props {
        request.Header.Set(prop.name, prop.value)
    }
}

func (httpReq *HttpReq) addCookies(request *http.Request) {
    for _, cookie := range httpReq.cookies {
        request.AddCookie(cookie)
    }
}

func (httpReq *HttpReq) postSettings(request *http.Request) {
    request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
}
