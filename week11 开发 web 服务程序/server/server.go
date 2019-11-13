package server

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

// 存储接收到的数据
type SubmitForm struct {
	Username string `form:"username" binding: "required"`
	Password string `form:"password" binding: "required"`
}

// 启动服务
func Start(port string) {
	m := martini.Classic()
	// 加入静态资源库
	m.Use(martini.Static("assets"))
	// 从 templates 路径中渲染 html 模板
	m.Use(render.Renderer())
	// 直接渲染 index.tmpl 模板
	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})
	// binding.Bind 将验证 SubmitForm 中的必填字段 username 和 password
	m.Post("/", binding.Bind(SubmitForm{}), func(sf SubmitForm, r render.Render) {
		// 获得数据
		info := SubmitForm{Username: sf.Username, Password: sf.Password}
		// 渲染 info 模板，将得到的用户信息传递回去
		r.HTML(200, "info", map[string]interface{}{"post": info})
	})
	// 运行
	m.RunOnAddr(":" + port)
}
