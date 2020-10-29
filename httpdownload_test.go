package gokits

import (
    "fmt"
    "github.com/stretchr/testify/assert"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    "time"
)

func TestTemp(t *testing.T) {
    reader := _HttpDownloadReader{}
    fmt.Println(reader.DelegateInterval)
}

func TestHttpDownload(t *testing.T) {
    a := assert.New(t)
    testServer := httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            ResponseText(w, "text")
        }))
    httpDownload := NewHttpDownload(WithRawURL(testServer.URL),
        WithDstDir("download"), WithFileName("test"), WithOverWrite(true))
    completed := false
    httpDownload.Start(
        WithDownloadDidFailWithError(func(download *HttpDownload, err error) {
            a.Fail("Should not fail")
            completed = true
        }),
        WithDownloadDidStarted(func(download *HttpDownload) {
            LOG.Debug("DownloadDidStarted")
        }),
        WithDownloadingWithProgress(func(download *HttpDownload, progress float64) {
            LOG.Debug("Downloading %.2f%%", progress)
        }),
        WithDownloadingProgressInterval(0),
        WithDownloadDidFinish(func(download *HttpDownload) {
            LOG.Debug("DownloadDidFinish")
            completed = true
        }))
    for !completed {
        time.Sleep(time.Second)
    }

    testFile, _ := os.Open("download/test")
    bytes, _ := ioutil.ReadAll(testFile)
    a.Equal("text", string(bytes))
    _ = testFile.Close()

    testServer = httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            ResponseText(w, "overwrite")
        }))
    httpDownload = NewHttpDownload(WithRawURL(testServer.URL),
        WithDstDir("download"), WithFileName("test"))
    completed = false
    httpDownload.Start(
        WithDownloadDidFailWithError(func(download *HttpDownload, err error) {
            a.Equal("DstFile exists, don't overwrite", err.Error())
            completed = true
        }),
        WithDownloadDidStarted(func(download *HttpDownload) {
            a.Fail("Should not download")
        }),
        WithDownloadingWithProgress(func(download *HttpDownload, progress float64) {
            a.Fail("Should not download")
        }),
        WithDownloadDidFinish(func(download *HttpDownload) {
            a.Fail("Should not download")
            completed = true
        }))
    for !completed {
        time.Sleep(time.Second)
    }

    testFile, _ = os.Open("download/test")
    bytes, _ = ioutil.ReadAll(testFile)
    a.Equal("text", string(bytes))
    _ = testFile.Close()
}
