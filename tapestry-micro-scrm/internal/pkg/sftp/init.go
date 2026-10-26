package sftp

import (
	"os"
	"time"

	"github.com/pkg/sftp"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/ctx"
	"gitlab-ce.k8s.tools.vchangyi.com/common/go-toolbox/log"
)

const (
	SSH_TIMEOUT = 10 * time.Second
)

type SFTP struct {
	Client *sftp.Client
}

func NewSFTPClient(sshConfig SSHConfig) (client *SFTP, err error) {
	sshClient, err := CreateSSHClient(sshConfig)
	if err != nil {
		return
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return
	}

	return &SFTP{
		Client: sftpClient,
	}, err
}

func (s *SFTP) Download(ctx ctx.Context, srcPath, dstPath string) {
	srcFile, _ := s.Client.Open(srcPath)
	dstFile, _ := os.Create(dstPath)
	defer func() {
		_ = srcFile.Close()
		_ = dstFile.Close()
	}()

	if _, err := srcFile.WriteTo(dstFile); err != nil {
		log.Info(ctx).Msgf("SFTP下载文件失败：%v; SrcFileName: %s", err, srcPath)
	}

	log.Info(ctx).Msgf("SFTP下载文件完成，SrcFileName: %s, DstFileName: %s", srcPath, dstPath)
}
