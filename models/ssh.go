package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
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

// create an ssh client with password
func SshClientWithPasswd(user, passwd, ip string, port int) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		Timeout:         SshTimeout,
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password(passwd)},
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), config)
	if err != nil {
		ourErr := fmt.Errorf("create ssh client fail: error: %w", err)
		beego.Error(ourErr)
		return nil, ourErr
	}
	return client, err
}

// create an ssh client with pem private key identity file
func SshClientWithPem(pemFilePath string, user string, ip string, port int) (*ssh.Client, error) {
	pemByte, err := ioutil.ReadFile(pemFilePath)
	if err != nil {
		ourErr := fmt.Errorf("read ssh private key file %s error: %w", pemFilePath, err)
		beego.Error(ourErr)
		return nil, ourErr
	}
	signer, err := ssh.ParsePrivateKey(pemByte)
	if err != nil {
		ourErr := fmt.Errorf("ssh.ParsePrivateKey error: %w", err)
		beego.Error(ourErr)
		return nil, ourErr
	}
	config := &ssh.ClientConfig{
		Timeout:         SshTimeout,
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), config)
	if err != nil {
		ourErr := fmt.Errorf("ssh.Dial error: %w", err)
		beego.Error(ourErr)
		return nil, ourErr
	}
	return client, nil
}
