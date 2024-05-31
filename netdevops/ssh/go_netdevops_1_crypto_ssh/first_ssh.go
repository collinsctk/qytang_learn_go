package main

import (
	"fmt"
	"go_netdevops_1_crypto_ssh/sshclient" // 替换为你的模块名
	"log"
)

func main() {
	username := "admin"
	password := "Cisc0123"
	ip := "10.10.1.1"
	commands := []string{
		"terminal length 0",   // 禁用分页
		"show running-config", // 显示配置
		"exit",                // 退出会话
	}

	result, err := sshclient.ExecuteSSHCommands(username, password, ip, commands)
	if err != nil {
		log.Fatalf("Error executing SSH commands: %s", err)
	}

	fmt.Println("Commands output:\n", result)
}
