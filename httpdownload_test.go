package gokits

import (
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    "time"
)

func TestHttpDownload(t *testing.T) {
    testServer := httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            ResponseText(w, "text")
        }))
    httpDownload := NewHttpDownload(WithRawURL(testServer.URL),
        WithDstDir("download"), WithFileName("test"), WithOverWrite(true))
    completed := false
    httpDownload.Start(
        WithDownloadDidFailWithError(func(download *HttpDownload, err error) {
            t.Errorf("Should not fail: %s", err.Error())
            completed = true
        }),
        WithDownloadDidStarted(func(download *HttpDownload) {
            LOG.Debug("DownloadDidStarted")
        }),
        WithDownloadingWithProgress(func(download *HttpDownload, progress float64) {
            LOG.Debug("Downloading %.2f%%", progress)
        }),
        WithDownloadDidFinish(func(download *HttpDownload) {
            LOG.Debug("DownloadDidFinish")
            completed = true
        }))
    for !completed {
        time.Sleep(time.Second)
    }

    testFile, _ := os.Open("download/test")
    bytes, _ := ioutil.ReadAll(testFile)
    if "text" != string(bytes) {
        t.Errorf("download/test file content should be \"text\"")
    }
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
            if "DstFile exists, don't overwrite" != err.Error() {
                t.Errorf("Should failed by overwrite")
            }
            completed = true
        }),
        WithDownloadDidStarted(func(download *HttpDownload) {
            t.Errorf("Should not download")
        }),
        WithDownloadingWithProgress(func(download *HttpDownload, progress float64) {
            t.Errorf("Should not download")
        }),
        WithDownloadDidFinish(func(download *HttpDownload) {
            t.Errorf("Should not download")
            completed = true
        }))
    for !completed {
        time.Sleep(time.Second)
    }

    testFile, _ = os.Open("download/test")
    bytes, _ = ioutil.ReadAll(testFile)
    if "text" != string(bytes) {
        t.Errorf("download/test file content should be \"text\"")
    }
    _ = testFile.Close()
}
