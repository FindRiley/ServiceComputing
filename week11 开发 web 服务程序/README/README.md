# Golang 开发简单 web 服务程序

## 要求

* 熟悉 go 服务器工作原理；
* 基于现有 web 库，编写一个简单 web 应用类似 cloudgo；
* 使用 curl 工具访问 web 程序；
* 对 web 执行压力测试。

## 实验过程

### 安装框架

```bash
go get github.com/go-martini/martini
go get github.com/martini-contrib/render
go get github.com/martini-contrib/binding
```

### 目录结构

>CloudGo
>
>* main.go
>
>- assets
>    - favicon.ico
>    - image.png
>- server
>    - server.go
>- templates
>    - index.tmpl
>    - info.tmpl

* 其中 assets 文件中放需要的静态资源，使用 `m.Use(martini.Static("assets"))` 加入该静态文件服务的文件夹。
* 需在 `templates` 目录下放入 `.tmpl` 模板文件。因为程序使用到 martini-render，默认情况下，`render.Renderer`中间件从 `templates` 目录中加载扩展名为 `.tmpl` 的模板。(见 [render-Loading Templates](https://github.com/martini-contrib/render#loading-templates))

### 写入模板文件

* 完成 index.tmpl 和 info.tmpl 网页模板的内容
* 在使用到上面 assets 目录中的资源时，比如 image.png。最后需要将 `<img src="../assets/image.png" alt="">` 改写为 `<img src="image.png" alt="">`，因为 `m.Use(martini.Static("assets"))`已经加载了资源目录，直接可以渲染。

### server.go

* 实现启动服务器，响应 Get 请求和 Post 请求。

* 使用 `Get()` 请求 index 页面

  ```go
  m.Get("/", func(r render.Render) {
  		r.HTML(200, "index", nil)
  	})
  ```

* 当在 index 页面中提交表单，将产生 Post 请求，请求 info 页面。

  ```go
  // 使用结构体存储接收到的数据
  type SubmitForm struct {
  	Username string `form:"username" binding: "required"`
  	Password string `form:"password" binding: "required"`
  }
  ```

  * 使用 binding 包，将原始请求映射到结构中。(见 [binding](https://github.com/martini-contrib/binding))

  ```go
  m.Post("/", binding.Bind(SubmitForm{}), func(sf SubmitForm, r render.Render) {
    // 获得数据
    info := SubmitForm{Username: sf.Username, Password: sf.Password}
    // 渲染 info 模板，将得到的用户信息传递回去
    r.HTML(200, "info", map[string]interface{}{"post": info})
  })
  ```

  * 这里将 `info` 结构体中的数据渲染到 `info.tmpl` 中的 `{{.post.*}}` 位置，其中的 `info.Username` 传到 `{{.post.Username}}` 位置， `info.Password` 传到 `{{.post.Password}}` 位置。

### 效果

![1573620075925](assets/1.png)

```bash
[shiyiloo@centosI CloudGo]$ ./main
[martini] listening on :8088 (development)
[martini] Started GET / for 127.0.0.1:46856
[martini] Completed 200 OK in 314.928µs
[martini] Started GET /favicon.ico for 127.0.0.1:46856
[martini] [Static] Serving /favicon.ico
[martini] Completed 200 OK in 26.896003ms
[martini] Started GET /favicon.ico for 127.0.0.1:46858
[martini] [Static] Serving /favicon.ico
[martini] Completed 200 OK in 163.666µs
```

![1573620212009](assets/2.png)

![1573620237199](assets/3.png)

```bash
[martini] Started POST / for 127.0.0.1:46862
[martini] Completed 200 OK in 421.322µs
```

### curl 测试

* `>` 表示请求，`<` 表示响应；
* `Content` 是请求的页面代码。

```bash
[shiyiloo@centosI ~]$ curl -v http://127.0.0.1:8088
* About to connect() to 127.0.0.1 port 8088 (#0)
*   Trying 127.0.0.1...
* Connected to 127.0.0.1 (127.0.0.1) port 8088 (#0)
> GET / HTTP/1.1
> User-Agent: curl/7.29.0
> Host: 127.0.0.1:8088
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: text/html; charset=UTF-8
< Date: Wed, 13 Nov 2019 04:54:49 GMT
< Content-Length: 836
< 
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Cloud Go</title>
    <link rel="Shortcut Icon" href="favicon.ico" type="image/x-icon" />
  </head>

  <body>
    <h1 style="color:#3E63A7" align="center">CloudGo</h1>
    <div id="main" width=100%>
    	<div id="left" style="width:32%;float:left">
        <p>Hi, Welcome to cloud go </p>
    		<p>Sign in here:</p>
    		
        <form method="post" action="/">
          <input type="text" name="username" placeholder="Username"><br /><br/>
          <input type="password" name="password" placeholder="Password"><br /><br />
          <input type="submit" value="sign in" id="submit">
        </form>
    	</div>
      
    	<div id="right" style="width:60%;float:right">
        <img src="image.png" alt="" width=92% />
      </div>
    </div>
  </body>

</html>
* Connection #0 to host 127.0.0.1 left intact

```

### ab 测试

* 安装工具，安装 Apache 会自动安装 ab ，单独安装如下

  ```bash
  sudo yum -y install httpd-tools
  ```

* 参数

  ```bash
  # 基本参数
  -n 执行的请求数量
  -c 并发请求个数
  # 其他参数
  -t 测试所进行的最大秒数
  -p 包含了需要 POST 的数据的文件
  -T POST 数据所使用的 Content-type 头信息
  -k 启用 HTTP KeepAlive 功能，即在一个HTTP会话中执行多个请求，默认时不启用
  ```

* 测试

  ```bash
  [shiyiloo@centosI ~]$ ab -n 5000 -c 100 http://127.0.0.1:8088/
  This is ApacheBench, Version 2.3 <$Revision: 1430300 $>
  Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
  Licensed to The Apache Software Foundation, http://www.apache.org/
  
  Benchmarking 127.0.0.1 (be patient)
  Completed 500 requests
  Completed 1000 requests
  Completed 1500 requests
  Completed 2000 requests
  Completed 2500 requests
  Completed 3000 requests
  Completed 3500 requests
  Completed 4000 requests
  Completed 4500 requests
  Completed 5000 requests
  Finished 5000 requests
  
  
  Server Software:        
  Server Hostname:        127.0.0.1
  Server Port:            8088
  
  Document Path:          /
  Document Length:        836 bytes
  
  Concurrency Level:      100
  Time taken for tests:   2.809 seconds
  Complete requests:      5000
  Failed requests:        0
  Write errors:           0
  Total transferred:      4765000 bytes
  HTML transferred:       4180000 bytes
  Requests per second:    1780.16 [#/sec] (mean)
  Time per request:       56.175 [ms] (mean)
  Time per request:       0.562 [ms] (mean, across all concurrent requests)
  Transfer rate:          1656.73 [Kbytes/sec] received
  
  Connection Times (ms)
                min  mean[+/-sd] median   max
  Connect:        0    3   4.4      2      42
  Processing:    10   53  19.0     51     147
  Waiting:        1   37  19.0     34     147
  Total:         11   56  19.2     54     149
  
  Percentage of the requests served within a certain time (ms)
    50%     54
    66%     61
    75%     66
    80%     71
    90%     82
    95%     89
    98%     93
    99%    108
   100%    149 (longest request)
  ```

  * Concurrency Level: 并发数
  * Time taken for tests: 完成所有请求所花时间
  * Complete requests: 总共完成的请求数
  * Failed requests: 失败的请求次数
  * Total transferred:  总共传输的字节数
  * HTML transferred:  html 页面传输的字节数
  * Requests per second:  每秒请求次数，吞吐率
  * Time per request (mean): 每次请求平均所用时间
  * Time per request (mean, across all concurrent requests): 服务器平均请求等待时间
  * Transfer rate: 传输速率

