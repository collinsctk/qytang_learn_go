package main

import (
	"fmt"
	"log"
	"time"

	"snmp_3_getall/influxdb_client" // 确保导入路径正确
	"snmp_3_getall/snmp_3_getall"
)

func main() {
	target := "10.10.1.11"
	community := "qytangro"

	// 初始化InfluxDB客户端
	influxClient, err := influxdb_client.NewInfluxDBClient("10.10.1.200", "8086", "qytdbuser", "Cisc0123", "qytdb")
	if err != nil {
		log.Fatalf("Error creating InfluxDB client: %v", err)
	}
	defer influxClient.Close()

	// 创建一个ticker，每10秒执行一次
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 获取统计数据
			stats, err := snmp_3_getall.SnmpGetAll(target, community)
			if err != nil {
				log.Printf("Error: %v", err)
				continue
			}

			fmt.Printf("%+v\n", stats)

			// 写入设备统计数据
			deviceTags := map[string]string{
				"device_ip":   stats.IPAddress,
				"device_type": "IOS-XE",
			}
			deviceFields := map[string]interface{}{
				"cpu_usage":   stats.CPU,
				"mem_usage":   stats.MemoryUsed,
				"mem_free":    stats.MemoryFree,
				"mem_percent": stats.MemoryPercent,
			}
			err = influxClient.WritePoints("router_monitor", deviceTags, deviceFields)
			if err != nil {
				log.Printf("Error writing device stats to InfluxDB: %v", err)
			}

			// 写入接口统计数据
			for _, iface := range stats.Interfaces {
				ifaceTags := map[string]string{
					"device_ip":      stats.IPAddress,
					"device_type":    "IOS-XE",
					"interface_name": iface.Name,
				}
				ifaceFields := map[string]interface{}{
					"status":    iface.Status,
					"in_bytes":  iface.InBytes,
					"out_bytes": iface.OutBytes,
				}
				err = influxClient.WritePoints("if_monitor", ifaceTags, ifaceFields)
				if err != nil {
					log.Printf("Error writing interface stats to InfluxDB: %v", err)
				}
			}

			fmt.Println("Data written to InfluxDB")
		}
	}
}
