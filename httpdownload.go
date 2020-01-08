package gokits

import (
    "errors"
    "io"
    "net/http"
    "net/url"
    "os"
    "path/filepath"
)

type HttpDownload struct {
    RawURL    string
    DstDir    string
    FileName  string
    OverWrite bool
}

type HttpDownloadOptions struct {
    RawURL    string
    DstDir    string
    FileName  string
    OverWrite bool
}

var defaultHttpDownloadOptions = HttpDownloadOptions{
    DstDir:    "",
    FileName:  "",
    OverWrite: false,
}

type HttpDownloadOption func(*HttpDownloadOptions)

func WithRawURL(rawURL string) HttpDownloadOption {
    return func(o *HttpDownloadOptions) { o.RawURL = rawURL }
}

func WithDstDir(dstDir string) HttpDownloadOption {
    return func(o *HttpDownloadOptions) { o.DstDir = dstDir }
}

func WithFileName(fileName string) HttpDownloadOption {
    return func(o *HttpDownloadOptions) { o.FileName = fileName }
}

func WithOverWrite(overWrite bool) HttpDownloadOption {
    return func(o *HttpDownloadOptions) { o.OverWrite = overWrite }
}

func NewHttpDownload(opts ...HttpDownloadOption) *HttpDownload {
    options := defaultHttpDownloadOptions
    for _, o := range opts {
        o(&options)
    }
    return &HttpDownload{
        RawURL:    options.RawURL,
        DstDir:    options.DstDir,
        FileName:  options.FileName,
        OverWrite: options.OverWrite,
    }
}

type HttpDownloadDelegate struct {
    DownloadDidFailWithError func(*HttpDownload, error)
    DownloadDidStarted       func(*HttpDownload)
    DownloadingWithProgress  func(*HttpDownload, float64)
    DownloadDidFinish        func(*HttpDownload)
}

type HttpDownloadDelegateOptions struct {
    DownloadDidFailWithError func(*HttpDownload, error)
    DownloadDidStarted       func(*HttpDownload)
    DownloadingWithProgress  func(*HttpDownload, float64)
    DownloadDidFinish        func(*HttpDownload)
}

var defaultHttpDownloadDelegateOptions = HttpDownloadDelegateOptions{}

type HttpDownloadDelegateOption func(*HttpDownloadDelegateOptions)

func WithDownloadDidFailWithError(f func(*HttpDownload, error)) HttpDownloadDelegateOption {
    return func(o *HttpDownloadDelegateOptions) { o.DownloadDidFailWithError = f }
}

func WithDownloadDidStarted(f func(*HttpDownload)) HttpDownloadDelegateOption {
    return func(o *HttpDownloadDelegateOptions) { o.DownloadDidStarted = f }
}

func WithDownloadingWithProgress(f func(*HttpDownload, float64)) HttpDownloadDelegateOption {
    return func(o *HttpDownloadDelegateOptions) { o.DownloadingWithProgress = f }
}

func WithDownloadDidFinish(f func(*HttpDownload)) HttpDownloadDelegateOption {
    return func(o *HttpDownloadDelegateOptions) { o.DownloadDidFinish = f }
}

func (download *HttpDownload) Start(opts ...HttpDownloadDelegateOption) {
    options := defaultHttpDownloadDelegateOptions
    for _, o := range opts {
        o(&options)
    }
    delegate := HttpDownloadDelegate{
        DownloadDidFailWithError: options.DownloadDidFailWithError,
        DownloadDidStarted:       options.DownloadDidStarted,
        DownloadingWithProgress:  options.DownloadingWithProgress,
        DownloadDidFinish:        options.DownloadDidFinish,
    }

    parsedURL, err := url.Parse(download.RawURL)
    if nil != err && nil != delegate.DownloadDidFailWithError {
        delegate.DownloadDidFailWithError(download, err)
        return
    }
    fileName := filepath.Base(parsedURL.Path)
    if 0 != len(download.FileName) {
        fileName = download.FileName
    }

    dstDir := "."
    if 0 != len(download.DstDir) {
        dstDir = download.DstDir
    }
    dir, err := os.Stat(dstDir)
    if nil == err && !dir.IsDir() && nil != delegate.DownloadDidFailWithError { // exists, but not directory
        delegate.DownloadDidFailWithError(download, errors.New("DstDir exists, but not directory"))
        return
    }
    err = os.MkdirAll(dstDir, 0777)
    if nil != err && nil != delegate.DownloadDidFailWithError {
        delegate.DownloadDidFailWithError(download, err)
        return
    }

    fileDir := PathJoin(dstDir, fileName)
    _, err = os.Stat(fileDir)
    if nil == err && !download.OverWrite && nil != delegate.DownloadDidFailWithError { // exists, don't overwrite
        delegate.DownloadDidFailWithError(download, errors.New("DstFile exists, don't overwrite"))
        return
    }

    go func() {
        file, err := os.Create(fileDir)
        if err != nil && nil != delegate.DownloadDidFailWithError {
            delegate.DownloadDidFailWithError(download, err)
            return
        }
        defer func() { _ = file.Close() }()

        res, err := http.Get(download.RawURL)
        if err != nil {
            delegate.DownloadDidFailWithError(download, err)
            return
        }
        defer func() { _ = res.Body.Close() }()

        reader := &_HttpDownloadReader{
            Reader: res.Body,
            Total:  res.ContentLength,
        }
        if nil != delegate.DownloadingWithProgress {
            reader.Delegate = &_HttpDownloadReaderDelegate{
                Downloading: func(reader *_HttpDownloadReader) {
                    delegate.DownloadingWithProgress(download, reader.Progress)
                },
            }
        }

        if nil != delegate.DownloadDidStarted {
            delegate.DownloadDidStarted(download)
        }
        _, _ = io.Copy(file, reader)
        if nil != delegate.DownloadDidFinish {
            delegate.DownloadDidFinish(download)
        }
    }()
}

type _HttpDownloadReaderDelegate struct {
    Downloading func(*_HttpDownloadReader)
}

type _HttpDownloadReader struct {
    io.Reader
    Total    int64
    Current  int64
    Progress float64
    Delegate *_HttpDownloadReaderDelegate
}

func (r *_HttpDownloadReader) Read(p []byte) (n int, err error) {
    n, err = r.Reader.Read(p)
    r.Current += int64(n)
    r.Progress = float64(r.Current*10000/r.Total) / 100
    if nil != r.Delegate && nil != r.Delegate.Downloading {
        r.Delegate.Downloading(r)
    }
    return
}
