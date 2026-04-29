package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/lib/upload"
	"os"
	"time"
)

type File struct {
}

// OssToken File
// @Tags File
// @Summary Get OSS token
// @Description Get OSS token
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/file/oss_token [get]
// @Security token
func (f *File) OssToken(c *gin.Context) {
	token := global.Oss.GetPolicyToken("")
	response.Success(c, token)
}

type FileBack struct {
	upload.CallbackBaseForm
	Url string `json:"url"`
}

// Notify Callback after successful upload
func (f *File) Notify(c *gin.Context) {

	res := global.Oss.Verify(c.Request)
	if !res {
		response.Fail(c, 101, response.TranslateMsg(c, "NoAccess"))
		return
	}
	fm := &FileBack{}
	if err := c.ShouldBind(fm); err != nil {
		fmt.Println(err)
	}
	fm.Url = global.Config.Oss.Host + "/" + fm.Filename
	response.Success(c, fm)

}

// Upload Upload file to local storage
// @Tags File
// @Summary Upload file to local storage
// @Description Upload file to local storage
// @Accept  multipart/form-data
// @Produce  json
// @Param file formData file true "File upload example"
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/file/upload [post]
// @Security token
func (f *File) Upload(c *gin.Context) {
	file, _ := c.FormFile("file")
	timePath := time.Now().Format("20060102") + "/"
	webPath := "/upload/" + timePath
	path := global.Config.Gin.ResourcesPath + webPath
	dst := path + file.Filename
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return
	}
	// upload file to the specified directory
	err = c.SaveUploadedFile(file, dst)
	if err != nil {
		return
	}
	// return the file's web path
	response.Success(c, gin.H{
		"url": webPath + file.Filename,
	})
}
