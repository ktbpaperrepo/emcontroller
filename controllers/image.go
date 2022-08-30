package controllers

import (
	"fmt"
	"github.com/astaxie/beego"

	"emcontroller/models"
)

type ImageController struct {
	beego.Controller
}

func (c *ImageController) Get() {
	images, err := models.ListRepoTags()
	if err != nil {
		beego.Error(fmt.Sprintf("error: %s", err.Error()))
		c.Data["imageList"] = []string{}
	}
	c.Data["imageList"] = images
	c.TplName = "image.tpl"
}

func (c *ImageController) Upload() {
	// get the file name
	f, fileHead, err := c.GetFile("imageFile")
	if err != nil {
		beego.Error(fmt.Sprintf("Open file error: %s", err.Error()))
		return
	}
	defer f.Close()
	fileName := fileHead.Filename
	beego.Info(fmt.Sprintf("filename: %s", fileName))

	// It seems that we do not need to save file to the server
	//// upload the file to server
	//var filePath string = models.UploadDir + fileName
	//err = c.SaveToFile("imageFile", filePath)
	//if err != nil {
	//	beego.Error(fmt.Sprintf("Upload file to server error: %s", err.Error()))
	//	c.Ctx.WriteString(fmt.Sprintf("Upload file to server error: %s", err.Error()))
	//	return
	//}
	//beego.Info(fmt.Sprintf("filename: %s, upload to the server successful.", fileName))

	// load the image file to the docker engine
	imageIdOrRepoTag, err := models.LoadImage(f)
	if err != nil {
		beego.Error(fmt.Sprintf("Load image error: %s", err.Error()))
		return
	}
	beego.Info(fmt.Sprintf("Load image to docker engine successfully, ID or RepoTag: %s", imageIdOrRepoTag))

	// add the tag to the image
	imageName := c.GetString("imageName")
	imageTag := c.GetString("imageTag")
	beego.Info(fmt.Sprintf("Add %s a new tag, name: %s, tag: %s", imageIdOrRepoTag, imageName, imageTag))
	repoTag, err := models.TagImage(imageIdOrRepoTag, imageName, imageTag)
	if err != nil {
		beego.Error(fmt.Sprintf("Tag image error: %s", err.Error()))
		return
	}

	// push the image to the Docker Registry
	resp, err := models.PushImage(repoTag)
	if err != nil {
		beego.Error(fmt.Sprintf("Push image error: %s", err.Error()))
		return
	}
	beego.Info(fmt.Sprintf("Push image successfully, resp: %s", resp))

	c.TplName = "uploadSuccess.tpl"
}
