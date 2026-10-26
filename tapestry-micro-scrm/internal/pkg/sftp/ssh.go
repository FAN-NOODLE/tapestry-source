package sftp

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
)

type SSHConfig struct {
	Host     string `json:"host"`
	Port     int64  `json:"port"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func CreateSSHClient(config SSHConfig) (*ssh.Client, error) {
	sshConfig := ssh.ClientConfig{
		User: config.UserName,
		Auth: []ssh.AuthMethod{ssh.Password(config.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: SSH_TIMEOUT,
	}

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

	return ssh.Dial("tcp", addr, &sshConfig)
}
