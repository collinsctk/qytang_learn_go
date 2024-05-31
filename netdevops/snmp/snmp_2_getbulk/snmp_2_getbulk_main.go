package main

import (
	"fmt"
	"log"
	"snmp_2_getbulk/snmp_2_getbulk_client" // 确保导入路径正确
)

func main() {
	target := "10.10.1.11"
	community := "qytangro"
	//oid := ".1.3.6.1.2.1.2.2.1.2" // 接口名称
	//oid := ".1.3.6.1.2.1.2.2.1.5" // 接口速率
	oid := ".1.3.6.1.2.1.2.2.1.10" // 接口输入字节数
	//oid := ".1.3.6.1.2.1.2.2.1.16" // 接口输出字节数
	nonRepeaters := uint8(0)
	maxRepetitions := uint32(10)

	result, err := snmp_2_getbulk_client.SnmpGetBulk(target, community, oid, nonRepeaters, maxRepetitions)
	if err != nil {
		log.Fatalf("SnmpGetBulk() err: %v", err)
	}

	// 使用for循环逐个打印每个接口描述
	for _, iface := range result {
		fmt.Println(iface)
	}
}
