package service

import (
	"gin-chat/models"
	"gin-chat/utils"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {

	url := upload(c.Writer, c.Request)

	if url != "" {
		utils.ResponseOK(c.Writer, url, "上传成功")
	}

}

func UploadAvatar(c *gin.Context) {
	url := upload(c.Writer, c.Request)

	if url != "" {

		userId, _ := strconv.Atoi(c.Request.FormValue("userid"))
		models.UpdateAvatarByUserId(uint(userId), url)
		utils.ResponseOK(c.Writer, url, "更新头像成功")
	}

}

func upload(writer gin.ResponseWriter, request *http.Request) string {
	// 实现文件上传功能
	file, head, err := request.FormFile("file")
	filetype := request.FormValue("filetype")

	if err != nil {
		utils.ResponseFail(writer, err.Error())
		return ""
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
		return ""
	}

	_, errr := io.Copy(desFile, file)
	if errr != nil {
		utils.ResponseFail(writer, errr.Error())
		return ""
	}

	return url

}
