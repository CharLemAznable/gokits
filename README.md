### gokits

[![Build Status](https://travis-ci.org/CharLemAznable/gokits.svg?branch=master)](https://travis-ci.org/CharLemAznable/gokits)
[![codecov](https://codecov.io/gh/CharLemAznable/gokits/branch/master/graph/badge.svg)](https://codecov.io/gh/CharLemAznable/gokits)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/CharLemAznable/gokits)
[![MIT Licence](https://badges.frapsoft.com/os/mit/mit.svg?v=103)](https://opensource.org/licenses/mit-license.php)
[![GoDoc](https://godoc.org/github.com/CharLemAznable/gokits?status.svg)](https://godoc.org/github.com/CharLemAznable/gokits)
![GitHub code size](https://img.shields.io/github/languages/code-size/CharLemAznable/gokits)

Go常用工具包.

##### boolean

  ```string to bool```工具.

##### cache

  重构来自[muesli/cache2go](https://github.com/muesli/cache2go/)

  在原有访问缓存过期策略的基础上, 增加写入缓存过期策略.

  在原有缓存Loader类型上, 增加异常返回值.

##### collections

  字符串数组/切片工具

##### cond

  条件表达式工具

##### gql

  数据库访问工具

##### hashset

  实现Java HashSet

##### hashtable

  实现Java HashTable
  
##### http

  HTTP工具

##### httpdownload

  HTTP文件下载

##### httpelf && httpserver

  HTTPServer工具

##### httpreq

  网络请求工具

##### json

  JSON序列化/反序列化工具

##### line_reader

  实现Java LineReader

##### log4go

  本地日志

  Fork from [https://github.com/alecthomas/log4go](https://github.com/alecthomas/log4go)

  Source Code from [http://code.google.com/p/log4go/](http://code.google.com/p/log4go/)

  Please see http://log4go.googlecode.com/

  Installation:
  - Run `goinstall log4go.googlecode.com/hg`

  Usage:
  - Add the following import:
  import l4g "log4go.googlecode.com/hg"

  Acknowledgements:
  - pomack
    For providing awesome patches to bring log4go up to the latest Go spec

##### path

  路径合并工具

##### properties

  实现Java Properties

##### random

  随机字符串
  
##### rsa

  RSA非对称加密算法工具, 包括密钥生成, 公钥加密/私钥解密, 私钥签名/公钥验签.

##### secret

  HmacSha256, AES工具

##### strconv

  字符串数字转换工具

##### yaml

  重构来自[kylelemons/go-gypsy](https://github.com/kylelemons/go-gypsy)

##### ycomb

  Go语言实现Y组合子

  参考:

  [使用Lambda 表达式编写递归三:实现 Y 组合子 - 鹤冲天 - 博客园](https://www.cnblogs.com/ldp615/archive/2013/04/10/recursive-lambda-expressions-3.html)
