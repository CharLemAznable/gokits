package gokits

import (
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

func TestNewProperties(t *testing.T) {
    a := assert.New(t)
    properties := NewProperties()
    a.NotNil(properties.mapper)
    a.NotNil(properties.defaults)
    a.Equal(0, properties.defaults.Size())

    file, err := os.Open("properties_test.properties")
    a.Nil(err)
    defer func() { _ = file.Close() }()

    err = properties.Load(file)
    a.Nil(err)

    properties.List(os.Stdout)

    a.Equal("127.0.0.1:6379", properties.GetProperty("redis.hosts"))
    a.Equal("127.0.0.1", properties.GetProperty("redis.host"))
    a.Equal("6379", properties.GetProperty("redis.port"))
    a.Equal("http://test.fshow.easy-hi.com/fshow-res", properties.GetProperty("static.prefix"))
    a.Equal("http://test.res.fshow.easy-hi.com/images/", properties.GetProperty("attach.imgPrefix"))
    a.Equal("http://test.fshow.easy-hi.com/fshow/", properties.GetProperty("root"))
    a.Equal("1", properties.GetProperty("static.version"))
    a.Equal("dev", properties.GetProperty("mode"))
    a.Equal("http://test.fshow.easy-hi.com:8000/fshow-res/dev/modules/", properties.GetProperty("template.path"))
    a.Equal("http://test.fshow.easy-hi.com:8000/fshow-res", properties.GetProperty("innerResPath.prefix"))
    a.Equal("http://test.res.fshow.easy-hi.com/musics/", properties.GetProperty("music.prefix"))
    a.Equal("http://127.0.0.1:8017/boss-biz/authorize/check-token", properties.GetProperty("origin.boss"))
    a.Equal("http://test.go.easy-hi.com/admin/scene/show/center/1508666666/initWxJs", properties.GetProperty("wxconfig"))
    a.Equal("nonExists", properties.GetPropertyDefault("wxconfig1", "nonExists"))

    filename := "properties_out.properties"
    writer, err := os.Create(filename)
    a.Nil(err)
    defer func() { _ = writer.Close() }()
    properties.Save(writer, "")

    file2, err := os.Open(filename)
    a.Nil(err)
    defer func() { _ = file2.Close(); _ = os.Remove(filename) }()
    properties2 := NewProperties()
    err = properties2.Load(file2)
    a.Nil(err)
    for _, name := range properties.StringPropertyNames() {
        a.Equal(properties.GetProperty(name), properties2.GetProperty(name))
    }
}
