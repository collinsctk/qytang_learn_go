package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"snmp_3_getall/influxdb_client"
	"snmp_3_getall/snmp_1_get_client"
	"snmp_3_getall/snmp_2_getbulk_client"
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

func main() {
	target := "10.10.1.11"
	community := "qytangro"

	client, err := influxdb_client.NewInfluxDBClient("10.10.1.200", "8086", "qytdbuser", "Cisc0123", "qytdb")
	if err != nil {
		log.Fatalf("Error creating InfluxDB client: %v", err)
	}
	defer client.Close()

	for {
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

		// 写入CPU和内存数据到InfluxDB
		cpuMemTags := map[string]string{
			"device_ip":   deviceInfo.IPAddress,
			"device_type": "IOS-XE",
		}
		cpuMemFields := map[string]interface{}{
			"cpu_usage": deviceInfo.CPU,
			"mem_usage": deviceInfo.MemoryUsed,
			"mem_free":  deviceInfo.MemoryFree,
		}
		err = client.WritePoints("router_monitor", cpuMemTags, cpuMemFields)
		if err != nil {
			log.Fatalf("Error writing CPU/Memory data to InfluxDB: %v", err)
		}

		// 写入接口数据到InfluxDB
		for _, iface := range deviceInfo.Interfaces {
			if iface.InBytes > 0 && iface.OutBytes > 0 {
				ifaceTags := map[string]string{
					"device_ip":      deviceInfo.IPAddress,
					"device_type":    "IOS-XE",
					"interface_name": iface.Name,
				}
				ifaceFields := map[string]interface{}{
					"in_bytes":  iface.InBytes,
					"out_bytes": iface.OutBytes,
				}
				err = client.WritePoints("if_monitor", ifaceTags, ifaceFields)
				if err != nil {
					log.Fatalf("Error writing interface data to InfluxDB: %v", err)
				}
			}
		}

		fmt.Printf("%+v\n", deviceInfo)
		time.Sleep(5 * time.Second)
	}
}
