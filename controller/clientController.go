package controller

import (
	"github.com/gin-gonic/gin"
	"go-bot/setting"
	"go-bot/utils/tools"
	"net/http"
)

func SavePartyPic(c *gin.Context) {
	code := 200
	data := make(map[string]string)
	msg := ""
	file, image, err := c.Request.FormFile("image")
	if err != nil {
		code = -1
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  "获取文件出错:" + err.Error(),
			"data": data,
		})
	}

	if image == nil {
		code = -1
	} else {
		imageName := tools.GetImageName(image.Filename)
		fullPath := tools.GetImageFullPath()
		savePath := tools.GetImagePath()

		src := fullPath + imageName
		if !tools.CheckImageExt(imageName) {
			code = -1
			msg = "文件格式不匹配"
			setting.Log.Warning("文件格式不匹配")
		} else if !tools.CheckImageSize(file) {
			code = -1
			msg = "文件太大"
			setting.Log.Warning("文件太大")
		} else {
			err := tools.CheckImage(fullPath)
			if err != nil {
				code = -1
				msg = "上传文件出错:" + err.Error()
				setting.Log.Warning(msg)
			} else if err := c.SaveUploadedFile(image, src); err != nil {
				code = -1
				msg = "上传文件出错:" + err.Error()
				setting.Log.Warning(msg)
			} else {
				data["image_url"] = tools.GetImageFullUrl(imageName)
				data["image_save_url"] = savePath + imageName
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func SavePartyInfo(c *gin.Context) error {

	return nil
}
