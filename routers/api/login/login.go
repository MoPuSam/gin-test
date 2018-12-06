package login

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/Unknwon/com"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/lunny/log"
	"net/http"

	"gin-test/models"
	"gin-test/pkg/e"
	"gin-test/pkg/setting"
	"gin-test/pkg/util"
)
const SecretKey = "oulam-develop"
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

//用户登录验证
func Login(c *gin.Context) {
	name:=c.PostForm("name")
	pwd:=c.PostForm("pwd")

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
	var data models.User
	if ! valid.HasErrors() {
		if models.ExistUserByName(name) {
			data = models.GetUserByName(name)
			if data.Password == password{
				//创建token
				token,_ :=util.CreateToken([]byte(SecretKey), "OULAM", data.ID, true)
				uid_cookie:=&http.Cookie{
					Name:   "token",
					Value:    token,
					Path:     "/",
					HttpOnly: false,
					MaxAge:   100,
				}
				http.SetCookie(c.Writer,uid_cookie)
				//将token写入数据表
				data.AuthToken = token
				models.EditUser(data.ID,data)
				code = e.SUCCESS
			} else {
				code = e.ERROR_PASSWORD
			}
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

