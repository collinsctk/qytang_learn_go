package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"snmp_3_getall/influxdb_client" // 确保导入路径正确
	"snmp_3_getall/snmp_3_getall"
)

type Device struct {
	Target    string
	Community string
}

var influxdb_ip = "10.10.1.200"
var influxdb_port = "8086"
var influxdb_user = "qytdbuser"
var influxdb_password = "Cisc0123"
var influxdb_db = "qytdb"
var loop_wait = 2 * time.Second
var max_retries = 3 // 最大重试次数

func main() {
	devices := []Device{
		{"10.10.1.11", "qytangro"},
		{"10.10.1.12", "qytangro"},
		{"10.10.1.13", "qytangro"},
		{"10.10.1.1", "qytangro"},
	}

	// 初始化InfluxDB客户端
	influxClient, err := influxdb_client.NewInfluxDBClient(influxdb_ip, influxdb_port, influxdb_user, influxdb_password, influxdb_db)
	if err != nil {
		log.Fatalf("Error creating InfluxDB client: %v", err)
	}
	defer influxClient.Close()

	// 创建一个ticker，每2秒执行一次
	ticker := time.NewTicker(loop_wait)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var wg sync.WaitGroup
			for _, device := range devices {
				wg.Add(1)
				go func(device Device) {
					defer wg.Done()
					queryAndWriteData(device, influxClient)
				}(device)
			}
			wg.Wait()
		}
	}
}

func queryAndWriteData(device Device, influxClient *influxdb_client.InfluxDBClient) {
	retries := 0
	for {
		stats, err := getSNMPData(device.Target, device.Community)
		if err != nil {
			log.Printf("Error: %v", err)
			if retries < max_retries {
				retries++
				time.Sleep(2 * time.Second) // 等待2秒后重试
				continue
			} else {
				log.Printf("Max retries reached for device: %v", device.Target)
				return
			}
		}

		//fmt.Printf("%+v\n", stats)

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
			if retries < max_retries {
				retries++
				time.Sleep(2 * time.Second) // 等待2秒后重试
				continue
			} else {
				log.Printf("Max retries reached for device: %v", device.Target)
				return
			}
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
				if retries < max_retries {
					retries++
					time.Sleep(2 * time.Second) // 等待2秒后重试
					continue
				} else {
					log.Printf("Max retries reached for device: %v", device.Target)
					return
				}
			}
		}

		fmt.Println("Data written to InfluxDB for device:", device.Target)
		break
	}
}

func getSNMPData(target, community string) (*snmp_3_getall.DeviceStats, error) {
	// 创建新的连接并获取统计数据
	stats, err := snmp_3_getall.SnmpGetAll(target, community)
	if err != nil {
		return nil, err
	}
	return stats, nil
}
