package models

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/astaxie/beego"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var cli *client.Client

func init() {
	cli = initDockerClient()
	beego.Info("Docker client initialized.")
}

func initDockerClient() *client.Client {
	var options []client.Opt
	options = append(options, client.FromEnv)
	options = append(options, client.WithAPIVersionNegotiation())

	dialer := &net.Dialer{
		Timeout: RequestTimeout,
	}
	options = append(options, client.WithDialContext(dialer.DialContext))

	options = append(options, client.WithHost("http://"+DockerEngine))
	options = append(options, client.WithTimeout(RequestTimeout))
	options = append(options, client.WithVersion("1.41"))

	c, err := client.NewClientWithOpts(options...)
	if err != nil {
		beego.Error(fmt.Sprintf("error: %s", err.Error()))
		panic(err)
	}
	return c
}

// list the images in the Docker Engine.
// curl http://192.168.100.36:19998/v1.41/images/json
func ListImages() ([]types.ImageSummary, error) {
	ctx := context.Background()
	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		beego.Error(fmt.Sprintf("error: %s", err.Error()))
	}
	return images, err
}

// load an image to the Docker Engine.
// curl -X POST http://192.168.100.36:19998/v1.41/images/load -H "Content-Type: application/x-tar" --data-binary '@/home/kubernetes/emcontroller/upload/helloworld2.tar'
func LoadImage(imageFile io.Reader) (string, error) {
	ctx := context.Background()
	resp, err := cli.ImageLoad(ctx, imageFile, false)
	if err != nil {
		beego.Error(fmt.Sprintf("error: %s", err.Error()))
		return "", err
	}
	defer resp.Body.Close()
	if !resp.JSON {
		beego.Error(fmt.Sprintf("Response is not Json"))
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		beego.Error(fmt.Sprintf("error: %s", err.Error()))
		return "", err
	}
	sb := string(b)
	beego.Info("Load image response: " + sb)

	// id: {"stream":"Loaded image ID: sha256:feb5d9fea6a5e9606aa995e879d862b825965ba48de054caab5ef356dc6b3412\n"}
	// tag: {"stream":"Loaded image: 192.168.100.36:5000/helloworld:latest\n"}
	var idOrRepoTag string
	if strings.Contains(sb, "Loaded image ID: sha256:") {
		tmp := sb
		tmp = strings.Split(tmp, "Loaded image ID: sha256:")[1]
		tmp = strings.Split(tmp, `\n"}`)[0]
		idOrRepoTag = tmp
	} else if strings.Contains(sb, "Loaded image: ") {
		tmp := sb
		tmp = strings.Split(tmp, "Loaded image: ")[1]
		tmp = strings.Split(tmp, `\n"}`)[0]
		idOrRepoTag = tmp
	}
	beego.Info("Image ID or RepoTag: " + idOrRepoTag)

	return idOrRepoTag, nil
}

// user should give a imageName and a targetTag to tag the image
// curl -v -X POST http://192.168.100.36:19998/v1.41/images/192.168.100.36:5000/ubuntu:latest/tag?repo=ubuntu1:vds123456
func TagImage(idOrRepoTag, imageName, targetTag string) (string, error) {
	ctx := context.Background()
	targetRepoTag := DockerRegistry + "/" + imageName + ":" + targetTag
	err := cli.ImageTag(ctx, idOrRepoTag, targetRepoTag)
	if err != nil {
		beego.Error(fmt.Sprintf("error: %s", err.Error()))
		return "", err
	}
	beego.Info(fmt.Sprintf("Successfully add %s a new RepoTag: %s", idOrRepoTag, targetRepoTag))
	return targetRepoTag, nil
}

// push the image from the Docker Engine to the Docker Registry
// curl  -v -X POST http://192.168.100.36:19998/v1.41/images/192.168.100.36:5000/helloworld:latest/push -H "X-Registry-Auth: eyJ1c2VybmFtZSI6InN0cmluZyIsInBhc3N3b3JkIjoic3RyaW5nIiwiZW1haWwiOiJzdHJpbmciLCJzZXJ2ZXJhZGRyZXNzIjoic3RyaW5nIn0K"
func PushImage(repoTag string) (string, error) {
	ctx := context.Background()
	respBody, err := cli.ImagePush(ctx, repoTag, types.ImagePushOptions{RegistryAuth: "eyJ1c2VybmFtZSI6InN0cmluZyIsInBhc3N3b3JkIjoic3RyaW5nIiwiZW1haWwiOiJzdHJpbmciLCJzZXJ2ZXJhZGRyZXNzIjoic3RyaW5nIn0K"})
	if err != nil {
		beego.Error(fmt.Sprintf("Push image error: %s", err.Error()))
		return "", err
	}
	b, err := io.ReadAll(respBody)
	if err != nil {
		beego.Error(fmt.Sprintf("Read response error: %s", err.Error()))
		return "", err
	}
	sb := string(b)
	beego.Info("Push image response: " + sb)
	return sb, nil
}
