package models

import (
	"testing"
)

func TestSshWithPem(t *testing.T) {
	sshClient, err := SshClientWithPem("conf/CLAAUDIAweifan.pem", SshUser, "10.92.1.198", SshPort)
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
