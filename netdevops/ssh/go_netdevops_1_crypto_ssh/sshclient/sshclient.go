package sshclient

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

// ExecuteSSHCommands 执行SSH命令
func ExecuteSSHCommands(username, password, ip string, commands []string) (string, error) {
	// 配置SSH客户端
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	// 建立SSH连接
	client, err := ssh.Dial("tcp", ip+":22", sshConfig)
	if err != nil {
		return "", fmt.Errorf("failed to dial: %w", err)
	}
	defer client.Close()

	// 创建一个新会话
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// 获取会话的标准输入和输出
	stdin, err := session.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("unable to setup stdin for session: %w", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("unable to setup stdout for session: %w", err)
	}

	// 启动一个伪终端
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return "", fmt.Errorf("request for pseudo terminal failed: %w", err)
	}

	// 启动shell
	if err := session.Shell(); err != nil {
		return "", fmt.Errorf("failed to start shell: %w", err)
	}

	// 创建一个缓冲区来存储输出
	var outputBuffer bytes.Buffer
	scanner := bufio.NewScanner(stdout)

	// 创建一个新的goroutine来读取输出
	go func() {
		for scanner.Scan() {
			outputBuffer.WriteString(scanner.Text() + "\n")
		}
	}()

	// 向设备发送命令
	for _, cmd := range commands {
		fmt.Fprintf(stdin, "%s\n", cmd)
		time.Sleep(1 * time.Second) // 等待命令执行完成
	}

	// 等待会话结束
	if err := session.Wait(); err != nil {
		return "", fmt.Errorf("session failed: %w", err)
	}

	return outputBuffer.String(), nil
}
