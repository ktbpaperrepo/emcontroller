package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"net/http"
)

type RegistryCatalog struct {
	Repositories []string `json:"repositories"`
}

type imageTags struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// get catalog from the docker registry.
// http://192.168.100.36:5000/v2/_catalog
func GetCatalog() ([]string, error) {
	url := "http://" + DockerRegistry + "/v2/_catalog"
	resp, err := http.Get(url)
	if err != nil {
		if err != nil {
			beego.Error(fmt.Sprintf("Get catalog error: %s", err.Error()))
			return []string{}, err
		}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		beego.Error(fmt.Sprintf("Read resp.Body error: %s", err.Error()))
		return []string{}, err
	}
	beego.Info("Get catalog response: " + string(b))
	var catalog RegistryCatalog
	err = json.Unmarshal(b, &catalog)
	if err != nil {
		beego.Error(fmt.Sprintf("Unmarshal error: %s", err.Error()))
		return []string{}, err
	}
	return catalog.Repositories, nil
}

// get tags of one image.
// curl http://192.168.100.36:5000/v2/helloworld12345/tags/list
func ListTags(imageName string) ([]string, error) {
	client := &http.Client{}
	url := fmt.Sprintf("http://%s/v2/%s/tags/list", DockerRegistry, imageName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		beego.Error(fmt.Sprintf("Create http request error: %s", err.Error()))
		return []string{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		beego.Error(fmt.Sprintf("Http request error: %s", err.Error()))
		return []string{}, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		beego.Error(fmt.Sprintf("Read resp.Body error: %s", err.Error()))
		return []string{}, err
	}
	beego.Info("get tags of image [" + imageName + "] response: " + string(b))
	var tags imageTags
	err = json.Unmarshal(b, &tags)
	if err != nil {
		beego.Error(fmt.Sprintf("Unmarshal error: %s", err.Error()))
		return []string{}, err
	}
	return tags.Tags, nil
}

// list all RepoTags in the Docker Registry
func ListRepoTags() ([]string, error) {
	repositories, err := GetCatalog()
	if err != nil {
		beego.Error(fmt.Sprintf("GetCatalog error: %s", err.Error()))
		return []string{}, err
	}
	var repoTags []string
	for _, repo := range repositories {
		tags, err := ListTags(repo)
		if err != nil {
			beego.Error(fmt.Sprintf("Repository %s, ListTags error: %s", repo, err.Error()))
			return []string{}, err
		}
		for _, tag := range tags {
			repoTags = append(repoTags, fmt.Sprintf("%s/%s:%s", DockerRegistry, repo, tag))
		}
	}
	return repoTags, nil
}
