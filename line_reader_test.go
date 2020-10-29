package gokits

import (
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

func TestNewLineReader(t *testing.T) {
    a := assert.New(t)
    file, err := os.Open("line_reader_test.properties")
    if err != nil {
        panic(err)
    }
    defer func() { _ = file.Close() }()

    lr := NewLineReader(file)
    lines := make([]string, 0)
    for limit, _ := lr.ReadLine(); limit >= 0; limit, _ = lr.ReadLine() {
        lines = append(lines, string(lr.lineBuf[:limit]))
    }

    a.Equal("redis.hosts=127.0.0.1:6379", lines[0])
    a.Equal("redis.host=127.0.0.1", lines[1])
    a.Equal("redis.port=6379", lines[2])
    a.Equal("static.prefix=http://test.fshow.easy-hi.com/fshow-res", lines[3])
    a.Equal("attach.imgPrefix=http://test.res.fshow.easy-hi.com/images/", lines[4])
    a.Equal("root=http://test.fshow.easy-hi.com/fshow/", lines[5])
    a.Equal("static.version=1", lines[6])
    a.Equal("mode=dev", lines[7])
    a.Equal("template.path=http://test.fshow.easy-hi.com:8000/fshow-res/dev/modules/", lines[8])
    a.Equal("innerResPath.prefix = http://test.fshow.easy-hi.com:8000/fshow-res", lines[9])
    a.Equal("music.prefix=http://test.res.fshow.easy-hi.com/musics/", lines[10])
    a.Equal("origin.boss=http://127.0.0.1:8017/boss-biz/authorize/check-token", lines[11])
    a.Equal("wxconfig=http://test.go.easy-hi.com/admin/scene/show/center/1508666666/initWxJs", lines[12])
}
