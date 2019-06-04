package gokits

import (
    "os"
    "testing"
)

func TestNewLineReader(t *testing.T) {
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

    if "redis.hosts=127.0.0.1:6379" != lines[0] {
        t.Fail()
    }
    if "redis.host=127.0.0.1" != lines[1] {
        t.Fail()
    }
    if "redis.port=6379" != lines[2] {
        t.Fail()
    }
    if "static.prefix=http://test.fshow.easy-hi.com/fshow-res" != lines[3] {
        t.Fail()
    }
    if "attach.imgPrefix=http://test.res.fshow.easy-hi.com/images/" != lines[4] {
        t.Fail()
    }
    if "root=http://test.fshow.easy-hi.com/fshow/" != lines[5] {
        t.Fail()
    }
    if "static.version=1" != lines[6] {
        t.Fail()
    }
    if "mode=dev" != lines[7] {
        t.Fail()
    }
    if "template.path=http://test.fshow.easy-hi.com:8000/fshow-res/dev/modules/" != lines[8] {
        t.Fail()
    }
    if "innerResPath.prefix = http://test.fshow.easy-hi.com:8000/fshow-res" != lines[9] {
        t.Fail()
    }
    if "music.prefix=http://test.res.fshow.easy-hi.com/musics/" != lines[10] {
        t.Fail()
    }
    if "origin.boss=http://127.0.0.1:8017/boss-biz/authorize/check-token" != lines[11] {
        t.Fail()
    }
    if "wxconfig=http://test.go.easy-hi.com/admin/scene/show/center/1508666666/initWxJs" != lines[12] {
        t.Fail()
    }
}
