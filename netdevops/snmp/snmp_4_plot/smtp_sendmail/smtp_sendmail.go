package smtp_sendmail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
)

// SendEmailWithAttachment 发送带有附件的电子邮件
func SendEmailWithAttachment(smtpHost, smtpPort, smtpUsername, smtpPassword, from, to, subject, body, attachmentPath string) error {
	// 读取附件文件
	attachmentData, err := os.ReadFile(attachmentPath)
	if err != nil {
		return fmt.Errorf("error reading attachment file: %v", err)
	}

	// 创建电子邮件
	e := email.NewEmail()
	e.From = from
	e.To = []string{to}
	e.Subject = subject
	e.Text = []byte(body)
	e.Attach(bytes.NewReader(attachmentData), "cpu_usage.png", "image/png")

	// 设置 SMTP 服务器连接
	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpHost)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	// 发送电子邮件
	err = e.SendWithTLS(fmt.Sprintf("%s:%s", smtpHost, smtpPort), auth, tlsConfig)
	if err != nil {
		return fmt.Errorf("error sending email: %v", err)
	}

	fmt.Println("Email sent successfully")
	return nil
}
