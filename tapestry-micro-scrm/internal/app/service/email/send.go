package email

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"micro-scrm/internal/app/http/request/email"
	"micro-scrm/internal/config"
	"micro-scrm/internal/pkg/helper"
	"os"
	"time"
)

type SendEmail struct {
	ServiceName string
}

func NewSendEmail() *SendEmail {
	return &SendEmail{
		ServiceName: "send_email",
	}
}

// Send 发送邮件
func (s *SendEmail) Send(sendMailReq emailReq.SendMailReq) error {
	// 获取url
	conf := config.CfgData.VchangyiEmail
	// 设置邮件内容
	mail := gomail.NewMessage()
	mail.SetHeader("From", conf.Username)
	mail.SetHeader("To", sendMailReq.ToAddress...)
	mail.SetHeader("Subject", sendMailReq.Title)
	mail.SetBody("text/html", sendMailReq.Content)
	if len(sendMailReq.FileList) > 0 {
		for _, fileInfo := range sendMailReq.FileList {
			if fileInfo.Url != "" {
				// 远程文件下载成本地文件
				tmpPath := fmt.Sprintf("./storage/email_attach/%s", fileInfo.Name)
				err := helper.DownloadFileToLocal(fileInfo.Url, tmpPath)
				if err != nil {
					return err
				}
				//本地文件
				mail.Attach(tmpPath) // 添加附件
			} else if fileInfo.LocalPath != "" {
				mail.Attach(fileInfo.LocalPath) // 添加附件
			}
		}

	}

	// 设置 SMTP 服务器信息
	d := gomail.NewDialer(conf.Host, conf.Port, conf.Username, conf.Password)

	// 发送邮件
	if err := d.DialAndSend(mail); err != nil {
		return err
	}
	return nil
}

// SendEmailByDiy 配置化发送邮件
func (s *SendEmail) SendEmailByDiy(sendMailReq emailReq.SendMailReq, emailConfig config.VchangyiEmailConfig) error {
	// 获取url
	conf := emailConfig
	// 设置邮件内容
	mail := gomail.NewMessage()
	mail.SetHeader("From", conf.Username)
	mail.SetHeader("To", sendMailReq.ToAddress...)
	mail.SetHeader("Subject", sendMailReq.Title)
	mail.SetBody("text/html", sendMailReq.Content)
	if len(sendMailReq.FileList) > 0 {
		for _, fileInfo := range sendMailReq.FileList {
			if fileInfo.Url != "" {
				// 远程文件下载成本地文件
				tmpPath := fmt.Sprintf("./storage/email_attach/%s", fileInfo.Name)
				err := helper.DownloadFileToLocal(fileInfo.Url, tmpPath)
				if err != nil {
					return err
				}
				//本地文件
				mail.Attach(tmpPath) // 添加附件
			} else if fileInfo.LocalPath != "" {
				mail.Attach(fileInfo.LocalPath) // 添加附件
			}
		}

	}

	// 设置 SMTP 服务器信息
	d := gomail.NewDialer(conf.Host, conf.Port, conf.Username, conf.Password)

	// 发送邮件
	if err := d.DialAndSend(mail); err != nil {
		return err
	}
	return nil
}

func (s *SendEmail) SendEmailByDiyGoFunc(req emailReq.SendMailReq, conf config.VchangyiEmailConfig) error {
	// 先下载所有附件到内存或临时文件
	attachments := make([]string, 0)

	// 处理URL附件
	for _, fileInfo := range req.FileList {
		if fileInfo.Url != "" {
			tmpPath := fmt.Sprintf("./storage/email_attach/%s", fileInfo.Name)
			if err := helper.DownloadFileToLocal(fileInfo.Url, tmpPath); err != nil {
				fmt.Printf("下载附件失败: %v\n", err)
				continue
			}
			attachments = append(attachments, tmpPath)
		}
	}

	// 发送完成后清理所有临时文件
	defer func() {
		for _, file := range attachments {
			os.Remove(file)
		}
	}()

	// 记录已有的本地文件路径
	localFiles := make([]string, 0)
	for _, fileInfo := range req.FileList {
		if fileInfo.LocalPath != "" {
			localFiles = append(localFiles, fileInfo.LocalPath)
		}
	}

	// 发送邮件
	for i, toAddress := range req.ToAddress {
		if i > 0 {
			time.Sleep(5 * time.Second)
		}

		mail := gomail.NewMessage()
		mail.SetHeader("From", conf.Username)
		mail.SetHeader("To", toAddress)
		mail.SetHeader("Subject", req.Title)
		mail.SetBody("text/html", req.Content)

		// 添加下载的临时文件
		for _, file := range attachments {
			mail.Attach(file)
		}

		// 添加原有的本地文件
		for _, file := range localFiles {
			mail.Attach(file)
		}

		d := gomail.NewDialer(conf.Host, conf.Port, conf.Username, conf.Password)

		if err := d.DialAndSend(mail); err != nil {
			fmt.Printf("发送失败到 %s: %v\n", toAddress, err)
		} else {
			fmt.Printf("成功发送到 %s\n", toAddress)
		}
	}

	fmt.Println("批量邮件发送完成")

	return nil
}
