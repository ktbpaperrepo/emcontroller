package models

import (
	"testing"
)

func TestSshWithPem(t *testing.T) {
	sshClient, err := SshClientWithPem("/root/.ssh/mc_id_rsa", SshRootUser, "10.234.234.157", SshPort)
	if err != nil {
		t.Errorf("SshClientWithPem error: %s", err)
	}
	defer sshClient.Close()
	output, err := SshOneCommand(sshClient, "pwd")
	if err != nil {
		t.Errorf("SshOneCommand error %s", err.Error())
	}
	t.Logf("output is: %s", string(output))
}

func TestSftpCopyFile(t *testing.T) {
	sshClient, err := SshClientWithPem("/root/.ssh/mc_id_rsa", SshRootUser, "10.234.234.100", SshPort)
	if err != nil {
		t.Errorf("SshClientWithPem error: %s", err.Error())
	}
	defer sshClient.Close()

	err = SftpCopyFile("/root/.kube/config", "/root/.kube/config", sshClient)
	if err != nil {
		t.Errorf("SftpCopyFile error: %s", err.Error())
	}
}
