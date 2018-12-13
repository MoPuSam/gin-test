package routers

import (
	"gin-test/middleware"
	"gin-test/routers/api/login"
	"gin-test/routers/api/wx"
	"github.com/gin-gonic/gin"
	"github.com/tommy351/gin-sessions"

	"gin-test/pkg/setting"
	"gin-test/routers/api/v1"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	//限制上传文件大小
	r.MaxMultipartMemory = 8 << 20

	//设置session midddleware
	store := sessions.NewCookieStore([]byte("mysession"))
	r.Use(sessions.Middleware("mysession", store))

	gin.SetMode(setting.RunMode)

	apiv0 := r.Group("/api")
	{
		apiv0.POST("/login", login.Login)
	}

	apiv1 := r.Group("/api/v1")
	apiv1.Use(middleware.AccessTokenMiddleware()) //必须声明在设置路由之前，否则无效
	{
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		//获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiv1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)

		//获取用户列表
		apiv1.GET("/users", v1.GetUsers)
		//获取指定用户
		apiv1.GET("/users/:id", v1.GetUser)
		//新增用户
		apiv1.POST("/users", v1.AddUser)
		//更新指定用户
		apiv1.PUT("/users/:id", v1.EditUser)
		//删除指定用户
		apiv1.DELETE("/users/:id", v1.DeleteUser)

		//上传文件
		apiv1.POST("/files", v1.Uploadfile)
	}
	apiv2 := r.Group("/api/v2")
	{
		apiv2.POST("/upload", v1.Uploadfile)
	}
	//微信服务器接口
	weChatR := r.Group("/weChatCore")
	//weChatR.Use(middleware.AccessTokenMiddleware())
	{
		weChatR.GET("/get", wx.GetWeChatCore)

	}
	return r
}
