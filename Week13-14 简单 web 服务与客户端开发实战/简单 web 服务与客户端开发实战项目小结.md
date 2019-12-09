# 项目小结

**作业要求** [简单 web 服务与客户端开发实战](https://pmlpml.github.io/ServiceComputingOnCloud/ex-services)

**项目地址** [SYSU-SimpleBlog](https://github.com/SYSU-SimpleBlog/)

**分工** 

本次作业分为前端开发和后端开发，前端 2 人，后端 4 人。我参与后端开发，和另一位同学负责后端有关 article 的三个 API 的设计，实现和测试。之后根据开发中的问题，修改完善了有关 comment 的功能代码和测试代码，并整合 user, article, comment 三部分的初步测试代码完成了对后端的初步测试。

## API 

* API 采用 REST v3 风格，设计了 user, article, comment 三种资源，共 6 个 API，如下

  ```json
  {
    "SignIn":"/simpleblog/user/signin",
    "DeleteArticleById":"/simpleblog/user/deleteArticle/{id}",
    "GetArticleById":"/simpleblog/user/article/{id}",
    "GetArticles":"/simpleblog/user/articles",
    "CreateComment":"/simpleblog/user/article/{id}/comment",
    "GetCommentsOfArticle":"/simpleblog/user/article/{id}/comments",
  }
  ```

* 使用 [swagger 编辑器](http://editor.swagger.io/)编写 API 文档，语法使用 yaml 语法。如下是 GetArticles 的 API 文档

  ```yaml
  paths:
    /user/articles:
      get:
        tags:
        - "article"
        summary: "Get articles"
        description: "Get the title of the request page of articles"
        operationId: "GetArticles"
        produces:
        - "application/json"
        parameters:
        - name: "page"
          in: "query"
          description: "the request page"
          required: true
          type: "string"
          x-exportParamName: "Page"
        responses:
          "200":
            description: "Successful Operation"
            schema:
              $ref: "#/definitions/ArticlesResponse"
          "404":
            description: "Not Found"
            schema:
              $ref: "#/definitions/inline_response_404"
  ```

## 实现

* [swagger 编辑器](http://editor.swagger.io/) 可以导出 go-Server ，本次后端开发是基于此框架实现的。

* 大体框架已经实现，主要需要实现的函数在三个以 api 为前缀的文件中，如下，分别对应 6 个 API。

  ```go
  func SignIn(w http.ResponseWriter, r *http.Request)
  func DeleteArticleById(w http.ResponseWriter, r *http.Request)
  func GetArticleById(w http.ResponseWriter, r *http.Request)
  func GetArticles(w http.ResponseWriter, r *http.Request)
  func CreateComment(w http.ResponseWriter, r *http.Request)
  func GetCommentsOfArticle(w http.ResponseWriter, r *http.Request)
  ```

**比如 getArticleById**

* 路由如下

  ```go
  Route{
  		"GetArticleById",
  		strings.ToUpper("Get"),
  		"/simpleblog/user/article/{id}",
  		GetArticleById,
  	},
  ```

* 响应过程大致为
  * 连接数据库
  * 解析 url，取到参数 Id
  * 打开 `article` 表查找是否有文章 id 为该 Id 的文章，找到则进行 json反序列化，然后返回该文章，否则返回相应的错误信息。

```go
func GetArticleById(w http.ResponseWriter, r *http.Request) {
	//connect to database
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	articleId := strings.Split(r.URL.Path, "/")[4]

	//	string to int
	Id, err := strconv.Atoi(articleId)
	fmt.Println(Id)
	if err != nil {
		reponse := ErrorResponse{"Wrong ArticleId"}
		JsonResponse(reponse, w, http.StatusBadRequest)
		return
	}

	var article Article
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Article"))
		if b != nil {
			v := b.Get(itob(Id))
			if v == nil {
				return errors.New("Article Not Exists")
			} else {
				_ = json.Unmarshal(v, &article)
				return nil
			}
		} else {
			return errors.New("Article Not Exists")
		}
	})

	if err != nil {
		response := ErrorResponse{err.Error()}
		JsonResponse(response, w, http.StatusNotFound)
		return
	}
	JsonResponse(article, w, http.StatusOK)
}
```

* 当然其他函数有更复杂的操作，在此不一一描述。

**问题**

* 测试的时候遇到了一个比较麻烦的问题，就是对带 token 的新建评论的测试，当时和另一个同学搜了很久都没搜到，最后通过一些零散的资料拼起来了测试命令，在此记录一下

  ```shell
  curl -H "Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzU2OTE4MzcsImlhdCI6MTU3NTY4ODIzN30.4Uw5rYZPCxB7uXNrVtn69Tmsy-831CgPnF8e555z-ko" http://localhost:8080/simpleblog/user/article/2/comment -X POST -d '{"content":"new content3","author":"user5"}'
  ```

## 测试

* 在此只做后端检测，除创建评论外，其他测试都能从浏览器可以直接访问

* 后端部分输出

  <left><img src="assets/1.png" width=60%></left>

* signin

```shell
http://localhost:8080/simpleblog/user/signin?username=user5&password=pass5
```

<left><img src="assets/2.png" width=100%></left>
```
"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzU5MDkyMTksImlhdCI6MTU3NTkwNTYxOX0.-cnC02_5I8bxqj7gib5gU1WhL4hSIR24DLi5sUCcfkw"
```

* getArticleById

```
http://localhost:8080/simpleblog/user/article/1
```

<left><img src="assets/3.png" width=100%></left>
- getArticles

```
http://localhost:8080/simpleblog/user/articles?page=1
```

 <left><img src="assets/4.png" width=80%></left>

- deleteArticle

```
http://localhost:8080/simpleblog/user/deleteArticle/1
```

  <left><img src="assets/5.png" width=100%></left>
- 再次 getArticles

<left><img src="assets/8.png" width=90%></left>

此时 id = 1 的文章已经被删除

- createComment 
  - 将之前登陆的 user 的 token 放在 Authorization 后

```
curl -H "Authorization:eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzU5MDkyMTksImlhdCI6MTU3NTkwNTYxOX0.-cnC02_5I8bxqj7gib5gU1WhL4hSIR24DLi5sUCcfkw" http://localhost:8080/simpleblog/user/article/2/comment -X POST -d '{"content":"new content3","author":"user5"}'
```

   <left><img src="assets/6.png" width=80%></left>

- getComments

```
http://localhost:8080/simpleblog/user/article/2/comments?page=1
```

   <left><img src="assets/7.png" width=70%></left>
可以看到第二条评论是刚刚加的



## 感想

本次作业是团队任务，团队合作可以激发人的责任心和团队精神，这是在个人作业中没有的，当身在一个团队中，不再是仅仅考虑自己累不累，更要考虑整个项目的进度。在此非常感谢和我一起为此次任务奋斗的组员。特别感谢与我一起开发 article 的另一位组员，在从医院看完病回来，忍着病痛完成了后端的最终测试，让我有最后的时间完成其他作业。