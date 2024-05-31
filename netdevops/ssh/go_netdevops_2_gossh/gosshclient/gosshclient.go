package sshclient

import (
	"fmt"
	"github.com/shenbowei/switch-ssh-go"
)

// ExecuteSSHCommands 执行SSH命令
func ExecuteSSHCommands(username, password, ipPort string, commands []string) (string, error) {
	// 禁用调试日志
	ssh.IsLogDebug = false

	// 获取设备品牌（供应商）
	brand, err := ssh.GetSSHBrand(username, password, ipPort)
	if err != nil {
		return "", fmt.Errorf("failed to get SSH brand: %w", err)
	}

	// 使用设备品牌执行命令
	result, err := ssh.RunCommandsWithBrand(username, password, ipPort, brand, commands...)
	if err != nil {
		return "", fmt.Errorf("failed to run commands: %w", err)
	}

	return result, nil
}
