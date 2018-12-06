package v1

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/lunny/log"
	"net/http"
	//"github.com/astaxie/beego/validation"
	"github.com/Unknwon/com"

	"gin-test/models"
	"gin-test/pkg/e"
	"gin-test/pkg/setting"
	"gin-test/pkg/util"
)

//获取多个用户信息
func GetUsers(c *gin.Context) {
	name := c.Query("name")

	maps := make(map[string]interface{})
	data := make(map[string]interface{})

	if name != "" {
		maps["name"] = name
	}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		maps["state"] = state
	}

	code := e.SUCCESS

	data["lists"] = models.GetTags(util.GetPage(c), setting.PageSize, maps)
	data["total"] = models.GetTagTotal(maps)

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : data,
	})
}

//新增用户
func AddUser(c *gin.Context) {
	//name := c.Query("name")//查询请求URL后面的参数
	//获取post里body的参数
	name:=c.PostForm("name")
	pwd:=c.PostForm("pwd")
	/*
	buf := make([]byte, 1024)
	n, _ := c.Request.Body.Read(buf)
	//fmt.Println(string(buf[0:n]))
	params :=string(buf[0:n])
	mappers:= strings.Split(params,"&")
	m := make(map[string]string, len(mappers))
	for _, r := range mappers {
		kv :=strings.Split(r,"=")
		m[kv[0]]=kv[1];
	}
	name:=m["name"]
	//fmt.Println(m["name"])
	pwd := m["pwd"]
	*/
	//fmt.Println(name+":"+pwd)
	ctx := md5.New()
	ctx.Write([]byte(pwd))
	password := hex.EncodeToString(ctx.Sum(nil))
	state := com.StrTo(c.DefaultQuery("state", "1")).MustInt()

	valid := validation.Validation{}
	valid.Required(name, "name").Message("名称不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	valid.Required(pwd, "pwd").Message("密码不能为空")
	valid.MaxSize(pwd, 12, "pwd").Message("名称最长为12字符")
	valid.MinSize(pwd, 6, "pwd").Message("名称最短为6字符")

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		if ! models.ExistUserByName(name) {
			code = e.SUCCESS
			models.AddUser(name, password,state)
		} else {
			code = e.ERROR_EXIST_USER
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]string),
	})
}
//修改用户信息
func EditUser(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}

			models.EditTag(id, data)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]string),
	})
}

//删除用户
func DeleteUser(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if ! valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			models.DeleteTag(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : make(map[string]string),
	})
}//获取单个用户
func GetUser(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	var data interface {}
	if ! valid.HasErrors() {
		if models.ExistArticleByID(id) {
			data = models.GetArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_USER
		}
	} else {
		for _, err := range valid.Errors {
			log.Info(err.Key, err.Message)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code" : code,
		"msg" : e.GetMsg(code),
		"data" : data,
	})
}
