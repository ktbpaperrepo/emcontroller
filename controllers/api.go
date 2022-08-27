package controllers

import (
	"crypto/md5"
	"fmt"
	"io"

	"github.com/astaxie/beego"

	"emcontroller/models"
)

type ApiController struct {
	beego.Controller
}

// http://localhost:20000/api/1231asf
func (c *ApiController) Get() {
	// get the value of dynamic router
	id := c.Ctx.Input.Param(":id")
	unix := 1599880013
	str := models.UnixToDate(unix)

	s1 := "2020-09-12 05:06:53"
	u := models.DateToUnix(s1)

	fmt.Println(models.GetUnix())
	fmt.Println(models.GetDate())

	h := md5.New()
	io.WriteString(h, "123456") //e10adc3949ba59abbe56e057f20f883e
	fmt.Printf("%x\n", h.Sum(nil))

	data := []byte("123456")
	fmt.Printf("%x\n", md5.Sum(data))

	fmt.Println(models.Md5("123456"))

	c.Ctx.WriteString("api interface---" + id + "---" + str + fmt.Sprintf("----%d", u))
}
