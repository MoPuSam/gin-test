package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Uploadfile(c *gin.Context) { //多文件上传

	form, _ := c.MultipartForm()
	files := form.File["upload[]"]

	for _, file := range files {
		log.Println(file.Filename)
		dstname := "G://upload/" + file.Filename
		fmt.Println("filename=", dstname)
		c.SaveUploadedFile(file, dstname)
	}
	c.String(http.StatusOK, "Uploaded...")

}
