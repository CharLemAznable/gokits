package gokits

import (
    "fmt"
    "os"
    "testing"
)

func TestNewProperties(t *testing.T) {
    properties := NewProperties()
    if nil == properties.mapper {
        t.Fail()
    }
    if nil != properties.defaults {
        t.Fail()
    }

    file, err := os.Open("properties_test.properties")
    if err != nil {
        t.Fail()
    }
    defer func() { _ = file.Close() }()

    if err = properties.Load(file); err != nil {
        t.Fail()
    }

    properties.List(os.Stdout)

    if properties.GetProperty("redis.hosts") != "127.0.0.1:6379" {
        t.Fail()
    }
    if properties.GetProperty("redis.host") != "127.0.0.1" {
        t.Fail()
    }
    if properties.GetProperty("redis.port") != "6379" {
        t.Fail()
    }
    if properties.GetProperty("static.prefix") != "http://test.fshow.easy-hi.com/fshow-res" {
        t.Fail()
    }
    if properties.GetProperty("attach.imgPrefix") != "http://test.res.fshow.easy-hi.com/images/" {
        t.Fail()
    }
    if properties.GetProperty("root") != "http://test.fshow.easy-hi.com/fshow/" {
        t.Fail()
    }
    if properties.GetProperty("static.version") != "1" {
        t.Fail()
    }
    if properties.GetProperty("mode") != "dev" {
        t.Fail()
    }
    if properties.GetProperty("template.path") != "http://test.fshow.easy-hi.com:8000/fshow-res/dev/modules/" {
        t.Fail()
    }
    if properties.GetProperty("innerResPath.prefix") != "http://test.fshow.easy-hi.com:8000/fshow-res" {
        t.Fail()
    }
    if properties.GetProperty("music.prefix") != "http://test.res.fshow.easy-hi.com/musics/" {
        t.Fail()
    }
    if properties.GetProperty("origin.boss") != "http://127.0.0.1:8017/boss-biz/authorize/check-token" {
        t.Fail()
    }
    if properties.GetProperty("wxconfig") != "http://test.go.easy-hi.com/admin/scene/show/center/1508666666/initWxJs" {
        t.Fail()
    }
    if properties.GetPropertyDefault("wxconfig1", "nonExists") != "nonExists" {
        t.Fail()
    }

    filename := "properties_out.properties"
    writer, err := os.Create(filename)
    if err != nil {
        t.Fail()
    }
    defer func() { _ = writer.Close() }()
    properties.Save(writer, "")

    file2, err := os.Open(filename)
    if err != nil {
        t.Fail()
    }
    defer func() { _ = file2.Close(); _ = os.Remove(filename) }()
    properties2 := NewProperties()
    if err = properties2.Load(file2); err != nil {
        t.Fail()
    }
    for _, name := range properties.StringPropertyNames() {
        if properties.GetProperty(name) != properties2.GetProperty(name) {
            fmt.Println(properties.GetProperty(name))
            fmt.Println(properties2.GetProperty(name))
            t.Fail()
        }
    }
}
