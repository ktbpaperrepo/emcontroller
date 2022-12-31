package models

import (
	"fmt"

	"github.com/astaxie/beego"
	"golang.org/x/crypto/ssh"
)

func SshOneCommand(client *ssh.Client, command string) ([]byte, error) {
	session, err := client.NewSession()
	if err != nil {
		beego.Error(fmt.Sprintf("Create ssh session fail: error: %s", err.Error()))
		return []byte{}, fmt.Errorf("Create ssh session for [%s] fail: error: %s", command, err.Error())
	}
	defer session.Close()

	output, err := session.Output(command)
	if err != nil {
		beego.Error(fmt.Sprintf("ssh execute [%s]: error: %s", command, err.Error()))
		beego.Error(fmt.Sprintf("ssh execute [%s]: output: %s", command, string(output)))
	} else {
		beego.Info(fmt.Sprintf("ssh execute [%s]: output: %s", command, string(output)))
	}
	return output, err
}
