package controllers

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"

	"emcontroller/models"
	"github.com/astaxie/beego"
)

type ImageController struct {
	beego.Controller
}

func (c *ImageController) Get() {
	repositories, err := models.GetCatalog()
	if err != nil {
		beego.Error(fmt.Sprintf("GetCatalog error: %s", err.Error()))
	}

	var repoTags map[string][]string = make(map[string][]string)
	for _, repo := range repositories {
		tags, err := models.ListTags(repo)
		if err != nil {
			beego.Error(fmt.Sprintf("Repository %s, ListTags error: %s", repo, err.Error()))
		}
		repoTags[repo] = tags
	}

	c.Data["dockerRegistry"] = models.DockerRegistry
	c.Data["imageList"] = repoTags
	c.TplName = "image.tpl"
}

// DeleteRepo delete a repository
func (c *ImageController) DeleteRepo() {
	repo := c.Ctx.Input.Param(":repo")

	beego.Info(fmt.Sprintf("Delete repository [%s]", repo))

	// use ssh to delete repository on docker registry
	dockerRegistryIP := beego.AppConfig.String("dockerRegistryIP")
	sshPort := 22
	sshUser := "root"
	sshPassword := beego.AppConfig.String("dockerRegiRootPasswd")

	config := &ssh.ClientConfig{
		Timeout:         5 * time.Second,
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password(sshPassword)},
	}

	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", dockerRegistryIP, sshPort), config)
	if err != nil {
		beego.Error(fmt.Sprintf("Create ssh client fail: error: %s", err.Error()))
	}
	defer sshClient.Close()

	// delete repository folder
	if _, err := models.SshOneCommand(sshClient, fmt.Sprintf("docker exec registry rm -rf /var/lib/registry/docker/registry/v2/repositories/%s", repo)); err != nil {
		beego.Error("ssh error: %s, exit", err.Error())
		return
	}

	// collect garbage
	if _, err := models.SshOneCommand(sshClient, "docker exec registry bin/registry garbage-collect /etc/docker/registry/config.yml"); err != nil {
		beego.Error("ssh error: %s, exit", err.Error())
		return
	}

	// restart docker
	if _, err := models.SshOneCommand(sshClient, "docker restart registry"); err != nil {
		beego.Error("ssh error: %s, exit", err.Error())
		return
	}

	beego.Info(fmt.Sprintf("Successful! Delete repository [%s]", repo))

	c.Ctx.ResponseWriter.WriteHeader(200)
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