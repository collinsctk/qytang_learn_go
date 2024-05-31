package main

import (
	"fmt"
	"log"
	"snmp_3_getall/snmp_1_get_client"
	"snmp_3_getall/snmp_2_getbulk_client"
	"strconv"
)

type Interface struct {
	Name     string `json:"name"`
	Status   bool   `json:"status"`
	InBytes  int    `json:"in_bytes"`
	OutBytes int    `json:"out_bytes"`
}

type DeviceInfo struct {
	CPU           int         `json:"cpu"`
	MemoryUsed    int         `json:"mem_used"`
	MemoryFree    int         `json:"mem_free"`
	MemoryPercent float64     `json:"mem_percent"`
	IPAddress     string      `json:"ip_address"`
	Interfaces    []Interface `json:"interface_list"`
}

// pprint 以漂亮的格式打印 DeviceInfo
func pprint(deviceInfo DeviceInfo) {
	fmt.Printf("{\n")
	fmt.Printf("  \"cpu\": %d,\n", deviceInfo.CPU)
	fmt.Printf("  \"mem_used\": %d,\n", deviceInfo.MemoryUsed)
	fmt.Printf("  \"mem_free\": %d,\n", deviceInfo.MemoryFree)
	fmt.Printf("  \"mem_percent\": %.2f,\n", deviceInfo.MemoryPercent)
	fmt.Printf("  \"ip_address\": \"%s\",\n", deviceInfo.IPAddress)
	fmt.Printf("  \"interface_list\": [\n")
	for i, iface := range deviceInfo.Interfaces {
		fmt.Printf("    {\n")
		fmt.Printf("      \"name\": \"%s\",\n", iface.Name)
		fmt.Printf("      \"status\": %t,\n", iface.Status)
		fmt.Printf("      \"in_bytes\": %d,\n", iface.InBytes)
		fmt.Printf("      \"out_bytes\": %d\n", iface.OutBytes)
		if i == len(deviceInfo.Interfaces)-1 {
			fmt.Printf("    }\n")
		} else {
			fmt.Printf("    },\n")
		}
	}
	fmt.Printf("  ]\n")
	fmt.Printf("}\n")
}

func main() {
	target := "10.10.1.11"
	community := "qytangro"

	deviceInfo := DeviceInfo{
		IPAddress: target,
	}

	// 获取CPU利用率
	cpu, err := snmp_1_get_client.SnmpGet(target, community, ".1.3.6.1.4.1.9.9.109.1.1.1.1.6.7")
	if err != nil {
		log.Fatalf("snmpGet() for CPU err: %v", err)
	}
	deviceInfo.CPU, _ = strconv.Atoi(cpu)

	// 获取内存使用情况
	memUsed, err := snmp_1_get_client.SnmpGet(target, community, ".1.3.6.1.4.1.9.9.109.1.1.1.1.12.7")
	if err != nil {
		log.Fatalf("snmpGet() for Memory Used err: %v", err)
	}
	deviceInfo.MemoryUsed, _ = strconv.Atoi(memUsed)

	memFree, err := snmp_1_get_client.SnmpGet(target, community, ".1.3.6.1.4.1.9.9.109.1.1.1.1.13.7")
	if err != nil {
		log.Fatalf("snmpGet() for Memory Free err: %v", err)
	}
	deviceInfo.MemoryFree, _ = strconv.Atoi(memFree)

	deviceInfo.MemoryPercent = float64(deviceInfo.MemoryUsed) / float64(deviceInfo.MemoryUsed+deviceInfo.MemoryFree) * 100

	// 获取接口信息
	ifNames, err := snmp_2_getbulk_client.SnmpGetBulk(target, community, ".1.3.6.1.2.1.2.2.1.2", 0, 10)
	if err != nil {
		log.Fatalf("SnmpGetBulk() for interface names err: %v", err)
	}

	ifStates, err := snmp_2_getbulk_client.SnmpGetBulk(target, community, ".1.3.6.1.2.1.2.2.1.8", 0, 10)
	if err != nil {
		log.Fatalf("SnmpGetBulk() for interface states err: %v", err)
	}

	ifInBytes, err := snmp_2_getbulk_client.SnmpGetBulk(target, community, ".1.3.6.1.2.1.2.2.1.10", 0, 10)
	if err != nil {
		log.Fatalf("SnmpGetBulk() for interface in bytes err: %v", err)
	}

	ifOutBytes, err := snmp_2_getbulk_client.SnmpGetBulk(target, community, ".1.3.6.1.2.1.2.2.1.16", 0, 10)
	if err != nil {
		log.Fatalf("SnmpGetBulk() for interface out bytes err: %v", err)
	}

	for i := range ifNames {
		ifName := ifNames[i]
		ifState, _ := strconv.Atoi(ifStates[i])
		ifIn, _ := strconv.Atoi(ifInBytes[i])
		ifOut, _ := strconv.Atoi(ifOutBytes[i])

		deviceInfo.Interfaces = append(deviceInfo.Interfaces, Interface{
			Name:     ifName,
			Status:   ifState == 1,
			InBytes:  ifIn,
			OutBytes: ifOut,
		})
	}

	pprint(deviceInfo)
}
