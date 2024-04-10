package service

import (
	"gin-chat/utils"
	"io"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	// 实现文件上传功能
	writer := c.Writer
	request := c.Request
	file, head, err := request.FormFile("file")
	filetype := request.FormValue("filetype")

	if err != nil {
		utils.ResponseFail(writer, err.Error())
		return
	}
	suffix := ".png"

	ofilName := strings.Split(head.Filename, ".")
	if len(ofilName) > 1 {
		suffix = "." + ofilName[len(ofilName)-1]
	}

	if filetype == ".mp3" {
		suffix = ".mp3"

	} else if filetype == ".mp4" {
		suffix = ".mp4"
	}

	fileName := utils.GetUUID() + suffix
	url := "./asset/upload/" + fileName
	desFile, err := os.Create(url)

	if err != nil {
		utils.ResponseFail(writer, err.Error())
		return
	}

	_, errr := io.Copy(desFile, file)
	if errr != nil {
		utils.ResponseFail(writer, errr.Error())
		return
	}

	utils.ResponseOK(writer, url, "上传成功")

}
