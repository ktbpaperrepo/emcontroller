package models

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/astaxie/beego"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func SshOneCommand(client *ssh.Client, command string) ([]byte, error) {
	beego.Info(fmt.Sprintf("Execute command on IP [%s]", client.Conn.RemoteAddr()))
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
		outErr := fmt.Errorf("create ssh client fail: error: %w", err)
		beego.Error(outErr)
		return nil, outErr
	}
	return client, err
}

// create an ssh client with pem private key identity file
func SshClientWithPem(pemFilePath string, user string, ip string, port int) (*ssh.Client, error) {
	pemByte, err := ioutil.ReadFile(pemFilePath)
	if err != nil {
		outErr := fmt.Errorf("read ssh private key file %s error: %w", pemFilePath, err)
		beego.Error(outErr)
		return nil, outErr
	}
	signer, err := ssh.ParsePrivateKey(pemByte)
	if err != nil {
		outErr := fmt.Errorf("ssh.ParsePrivateKey error: %w", err)
		beego.Error(outErr)
		return nil, outErr
	}
	config := &ssh.ClientConfig{
		Timeout:         SshTimeout,
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", ip, port), config)
	if err != nil {
		outErr := fmt.Errorf("ssh.Dial error: %w", err)
		beego.Error(outErr)
		return nil, outErr
	}
	return client, nil
}

func SftpCopyFile(srcPath, dstPath string, sshClient *ssh.Client) error {
	beego.Info(fmt.Sprintf("SFTP copy file [local:%s] to [%s:%s].", srcPath, sshClient.Conn.RemoteAddr(), dstPath))

	// open an SFTP session over an existing ssh connection.
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		if err != nil {
			outErr := fmt.Errorf("create SFTP client, error: %w", err)
			beego.Error(outErr)
			return outErr
		}
	}
	defer sftpClient.Close()

	// Open the source file
	srcFile, err := os.Open(srcPath)
	if err != nil {
		outErr := fmt.Errorf("open source file %s, error: %w", srcPath, err)
		beego.Error(outErr)
		return outErr
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := sftpClient.Create(dstPath)
	if err != nil {
		outErr := fmt.Errorf("create the destination file %s:%s, error: %w", sshClient.Conn.RemoteAddr(), dstPath, err)
		beego.Error(outErr)
		return outErr
	}
	defer dstFile.Close()

	// write from source file to destination file
	if _, err := dstFile.ReadFrom(srcFile); err != nil {
		outErr := fmt.Errorf("write from source file %s to the destination file %s:%s, error: %w", srcPath, sshClient.Conn.RemoteAddr(), dstPath, err)
		beego.Error(outErr)
		return outErr
	}

	return nil
}
