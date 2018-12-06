package middleware

import (
	"fmt"
	"gin-test/models"
	"gin-test/pkg/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)
//秘钥
const SecretKey = "oulam-develop"
func AccessTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authtoken, _ :=c.Cookie("token")
		if authtoken==""{
			// 没有提供权限token
			c.JSON(401, gin.H{
				"status": 401,
				"msg":    "cookie缺少authtoken！",
			})
			//终止请求，直接返回
			c.Abort()
			return
		}else {
			//auth, err := DecodeAuthV1(authtoken)
			auth, err := util.ParseToken(authtoken, []byte(SecretKey))
			if err != nil {
				// 权限信息不完整
				c.JSON(401,gin.H{
					"status":401,
					"msg":"权限信息不完整！",
				})
				c.Abort()
				return
			}
			//account := &models.User{}
			//o := common.NewOrm()
			//o.QueryTable("account").Filter("id", auth["id"]).One(account, "id", "account_token", "auth_token")
			var account models.User
			fmt.Println("middleware:",auth.(jwt.MapClaims)["uid"])
			//auth.(jwt.MapClaims)["uid"]返回的值一开始为float64，若直接转uint，运行时报错
			id:=auth.(jwt.MapClaims)["uid"].(float64)
			fmt.Println("middleware-unit:",uint(id))
			account = models.GetUserById(uint(id))
			if account.ID == 0 || account.AuthToken != authtoken {
				// 权限信息伪造或者已经失效
				c.JSON(401, gin.H{
					"status": 401,
					"msg":   "权限信息已经失效,请重新登录",
				})
				c.Abort()
				return
			}

		}
		c.Next()
	}
}
/**
func NewAccountSsdbCache(userid int64,val string)(err error){
	ssdb,err:=common.NewSsdbClient()
	if err!=nil {
		return err
	}
	defer ssdb.Close()
	if err:=ssdb.Set(MakeSSdbCacheKey(userid),val,60);err!=nil{
		return err
	}
	return nil
}
func NewAccountSsdbCache(userid int64,val string)(err error){
	ssdb,err:=common.NewSsdbClient()
	if err!=nil {
		return err
	}
	defer ssdb.Close()
	if err:=ssdb.Set(MakeSSdbCacheKey(userid),val,60);err!=nil{
		return err
	}
	return nil
}
**/