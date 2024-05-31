package main

import (
	"fmt"
	"log"
	"snmp_1_get/snmp_1_get_client" // 导入路径根据模块名和包名调整
)

func main() {
	target := "10.10.1.11"
	community := "qytangro"
	oid := ".1.3.6.1.2.1.1.5.0" // 主机名
	//oid := ".1.3.6.1.4.1.9.9.109.1.1.1.1.6.7"  // CPU利用率
	//oid := ".1.3.6.1.4.1.9.9.109.1.1.1.1.12.7" // 内存使用
	//oid := ".1.3.6.1.4.1.9.9.109.1.1.1.1.13.7" // 内存空闲

	result, err := snmp_1_get_client.SnmpGet(target, community, oid)
	if err != nil {
		log.Fatalf("snmpGet() err: %v", err)
	}

	fmt.Println(result)
}
