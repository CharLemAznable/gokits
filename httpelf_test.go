package gokits

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestDumpRequest(t *testing.T) {
    testServer := httptest.NewServer(DumpRequest(
        func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
        }))
    _, err := NewHttpReq(testServer.URL).Get()
    if nil != err {
        t.Errorf("Should has no error")
    }
}

func TestGzipHandlerFunc(t *testing.T) {
    testServer := httptest.NewServer(GzipResponse(
        func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
        }))
    _, err := NewHttpReq(testServer.URL).Get()
    if nil != err {
        t.Errorf("Should has no error")
    }
}

func TestServeModelContext(t *testing.T) {
    testServer := httptest.NewServer(ServeModelContext(
        func(w http.ResponseWriter, r *http.Request) {
            _, ok := r.Context().(*ModelCtx)
            if !ok {
                t.Errorf("Type should be ModelCtx")
            }
            w.WriteHeader(http.StatusOK)
        }))
    _, err := NewHttpReq(testServer.URL).Get()
    if nil != err {
        t.Errorf("Should has no error")
    }
}

func TestServeModelContextWithValueFunc(t *testing.T) {
    testServer := httptest.NewServer(ServeModelContextWithValueFunc(
        func(w http.ResponseWriter, r *http.Request) {
            modelCtx, ok := r.Context().(*ModelCtx)
            if !ok {
                t.Errorf("Type should be ModelCtx")
            }
            if modelCtx.Value("key") != "value" {
                t.Errorf("Context[key] should be \"value\"")
            }
            w.WriteHeader(http.StatusOK)
        }, func() (string, interface{}) {
            return "key", "value"
        }))
    _, err := NewHttpReq(testServer.URL).Get()
    if nil != err {
        t.Errorf("Should has no error")
    }
}

func TestServeMethod(t *testing.T) {
    testServer := httptest.NewServer(ServeGet(
        func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
        }))
    code, _, _ := NewHttpReq(testServer.URL).testGet()
    if code != http.StatusOK {
        t.Errorf("Should response http.StatusOK")
    }
    code, _, _ = NewHttpReq(testServer.URL).testPost()
    if code != http.StatusNotFound {
        t.Errorf("Should response http.StatusNotFound")
    }

    testServer = httptest.NewServer(ServePost(
        func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
        }))
    code, _, _ = NewHttpReq(testServer.URL).testGet()
    if code != http.StatusNotFound {
        t.Errorf("Should response http.StatusNotFound")
    }
    code, _, _ = NewHttpReq(testServer.URL).testPost()
    if code != http.StatusOK {
        t.Errorf("Should response http.StatusOK")
    }
}

func TestServeAjax(t *testing.T) {
    testServer := httptest.NewServer(ServeAjax(
        func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusOK)
        }))
    code, _, _ := NewHttpReq(testServer.URL).testGet()
    if code != http.StatusNotFound {
        t.Errorf("Should response http.StatusNotFound")
    }
}

func TestMinify(t *testing.T) {
    html := MinifyHTML(`<!DOCTYPE html><html>
<head><meta charset="utf-8"/><script src="http://code.jquery.com/jquery-latest.min.js"></script></head>
<body>
<div id="wrap"><div id="header"><h1>html在线工具</h1>
<!--   如果有用，请别忘了推荐给你的朋友：		-->
<!--   Html在线美化、格式化：https://tool.lu/html   -->
</div>
<div id="main">
	<!-- [history] -->
	<dl>
	<dt>v1.0</dt> <dd>2011-06-05 Html工具上线</dd>
	<dt>v1.1</dt> <dd>2012-01-14 修复美化功能，增加压缩</dd>
	<dt>v1.2</dt> <dd>2012-07-20 增加清除链接功能</dd>
	<dt>v1.3</dt> <dd>2014-08-05 修改 html 压缩引擎</dd>
	<dt>v1.4</dt> <dd>2014-08-09 增加转换为js变量的功能</dd>
	</dl>
</div>
<div id="footer">This is just an example.</div>
</div>
</body></html>`, false)
    expectHtml := "<!doctype html><meta charset=utf-8><script src=http://code.jquery.com/jquery-latest.min.js></script><div id=wrap><div id=header><h1>html在线工具</h1></div><div id=main><dl><dt>v1.0<dd>2011-06-05 Html工具上线<dt>v1.1<dd>2012-01-14 修复美化功能，增加压缩<dt>v1.2<dd>2012-07-20 增加清除链接功能<dt>v1.3<dd>2014-08-05 修改 html 压缩引擎<dt>v1.4<dd>2014-08-09 增加转换为js变量的功能</dl></div><div id=footer>This is just an example.</div></div>"
    if expectHtml != html {
        t.Errorf("Should minified")
    }

    css := MinifyCSS(`/*   美化：格式化代码，使之容易阅读			*/
/*   净化：将代码单行化，并去除注释   */
/*   整理：按照一定的顺序，重新排列css的属性   */
/*   优化：将css的长属性值优化为简写的形式   */
/*   压缩：将代码最小化，加快加载速度   */

/*   如果有用，请别忘了推荐给你的朋友：		*/
/*   css在线美化、压缩：https://tool.lu/css   */
/*   v1.1 2012-05-11   */
/*   v1.2 2015-04-30   */
/*   v1.3 2015-06-01 修复 css 压缩的 bug  */
/*   v1.4 2015-07-31 增加 css 优化 功能  */
/*   v1.5 2016-06-18 增加 px转rem 功能  */
/*   v1.6 2017-08-03 增加 加范围功能  */
/*   v1.7 2018-12-30 增加 px转rpx 功能  */

.css3 {
	box-shadow: 0 0;
	width: calc(100% + 2em);
	font-size: 24px;
}

/*   以下是演示代码				*/

body, div, dl, dt, dd, ul, ol, li,
h1, h2, h3, h4, h5, h6, pre, code,
form, fieldset, legend, input, button,
textarea, p, blockquote, th, td {
    margin: 0;
    padding: 0;
}
fieldset, img {
    border: 0;
}
/* remember to define focus styles! */
:focus {
    outline: 0;
}
address, ctoolion, cite, code, dfn,
em, strong, th, var, optgroup {
    font-style: normal;
    font-weight: normal;
}
 
h1, h2, h3, h4, h5, h6 {
    font-size: 100%;
    font-weight: normal;
}
abbr, acronym {
    border: 0;
    font-variant: normal;
}
 
input, button, textarea,
select, optgroup, option {
    font-family: inherit;
    font-size: inherit;
    font-style: inherit;
    font-weight: inherit;
}
code, kbd, samp, tt {
    font-size: 100%;
}
/*@purpose To enable resizing for IE */
/*@branch For IE6-Win, IE7-Win */
input, button, textarea, select {
    *font-size: 100%;
}
body {
    line-height: 1.5;
}
ol, ul {
    list-style: none;
}
/* tables still need 'cellspacing="0"' in the markup */
table {
    border-collapse: collapse;
    border-spacing: 0;
}
ctoolion, th {
    text-align: left;
}
sup, sub {
    font-size: 100%;
    vertical-align: baseline;
}
/* remember to highlight anchors and inserts somehow! */
:link, :visited , ins {
    text-decoration: none;
}
blockquote, q {
    quotes: none;
}
blockquote:before, blockquote:after,
q:before, q:after {
    content: '';
    content: none;
}`, false)
    expectCss := ".css3{box-shadow:0 0;width:calc(100% + 2em);font-size:24px}body,div,dl,dt,dd,ul,ol,li,h1,h2,h3,h4,h5,h6,pre,code,form,fieldset,legend,input,button,textarea,p,blockquote,th,td{margin:0;padding:0}fieldset,img{border:0}:focus{outline:0}address,ctoolion,cite,code,dfn,em,strong,th,var,optgroup{font-style:normal;font-weight:400}h1,h2,h3,h4,h5,h6{font-size:100%;font-weight:400}abbr,acronym{border:0;font-variant:normal}input,button,textarea,select,optgroup,option{font-family:inherit;font-size:inherit;font-style:inherit;font-weight:inherit}code,kbd,samp,tt{font-size:100%}input,button,textarea,select{*font-size:100%}body{line-height:1.5}ol,ul{list-style:none}table{border-collapse:collapse;border-spacing:0}ctoolion,th{text-align:left}sup,sub{font-size:100%;vertical-align:baseline}:link,:visited,ins{text-decoration:none}blockquote,q{quotes:none}blockquote:before,blockquote:after,q:before,q:after{content:'';content:none}"
    if expectCss != css {
        t.Errorf("Should minified")
    }

    js := MinifyJs(`/*   美化：格式化代码，使之容易阅读			*/
/*   净化：去掉代码中多余的注释、换行、空格等	*/
/*   压缩：将代码压缩为更小体积，便于传输		*/
/*   解压：将压缩后的代码转换为人可以阅读的格式	*/
/*   混淆：将代码的中变量名简短化以减小体积，但可读性差，经混淆后的代码无法还原	*/

/*   如果有用，请别忘了推荐给你的朋友：		*/
/*   javascript在线美化、净化、压缩、解压：https://tool.lu/js   */

/*   以下是演示代码				*/
var Inote = {};
Inote.JSTool = function(options) {
this.options = options || {};
};
Inote.JSTool.prototype = {
	_name: 'Javascript工具',
_history: {
		'v1.0': ['2011-01-18', 'javascript工具上线'],
		'v1.1': ['2012-03-23', '增加混淆功能'],
		'v1.2':	['2012-07-21', '升级美化功能引擎'],
		'v1.3': ['2014-03-01', '升级解密功能，支持eval,window.eval,window["eval"]等的解密'],
		'v1.4':	['2014-08-05', '升级混淆功能引擎'],
		'v1.5':	['2014-08-09', '升级js压缩引擎'],
		'v1.6':	['2015-04-11', '升级js混淆引擎'],
		'v1.7':	['2017-02-12', '升级js混淆引擎']
	},
	options: {},
	getName: function() {return this._name;},
	getHistory: function() {
		return this._history;}
};
var jstool = new Inote.JSTool();`, false)
    expectJs := "var Inote={};Inote.JSTool=function(options){this.options=options||{};};Inote.JSTool.prototype={_name:'Javascript工具',_history:{'v1.0':['2011-01-18','javascript工具上线'],'v1.1':['2012-03-23','增加混淆功能'],'v1.2':['2012-07-21','升级美化功能引擎'],'v1.3':['2014-03-01','升级解密功能，支持eval,window.eval,window[\"eval\"]等的解密'],'v1.4':['2014-08-05','升级混淆功能引擎'],'v1.5':['2014-08-09','升级js压缩引擎'],'v1.6':['2015-04-11','升级js混淆引擎'],'v1.7':['2017-02-12','升级js混淆引擎']},options:{},getName:function(){return this._name;},getHistory:function(){return this._history;}};var jstool=new Inote.JSTool();"
    if expectJs != js {
        t.Errorf("Should minified but %s", js)
    }
}

func TestEmptyHandler(t *testing.T) {
    testServer := httptest.NewServer(EmptyHandler)
    code, result, err := NewHttpReq(testServer.URL).testGet()
    if code != http.StatusOK {
        t.Errorf("Should response http.StatusOK")
    }
    if result != "" {
        t.Errorf("Should response empty string")
    }
    if err != nil {
        t.Errorf("Should response no error")
    }
}
