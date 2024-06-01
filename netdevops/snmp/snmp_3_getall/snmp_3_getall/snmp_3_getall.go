package snmp_3_getall

import (
	"snmp_3_getall/snmp_1_get_client"
	"snmp_3_getall/snmp_2_getbulk_client"
	"strconv"
)

type InterfaceStats struct {
	Name     string
	Status   bool
	InBytes  int
	OutBytes int
}

type DeviceStats struct {
	CPU           int
	MemoryUsed    int
	MemoryFree    int
	MemoryPercent float64
	IPAddress     string
	Interfaces    []InterfaceStats
}

func SnmpGetAll(target, community string) (*DeviceStats, error) {
	cpuOID := ".1.3.6.1.4.1.9.9.109.1.1.1.1.6.7"
	memUsedOID := ".1.3.6.1.4.1.9.9.109.1.1.1.1.12.7"
	memFreeOID := ".1.3.6.1.4.1.9.9.109.1.1.1.1.13.7"

	// 获取CPU信息
	cpuStr, err := snmp_1_get_client.SnmpGet(target, community, cpuOID)
	if err != nil {
		return nil, err
	}
	cpu, _ := strconv.Atoi(cpuStr)

	// 获取内存使用信息
	memUsedStr, err := snmp_1_get_client.SnmpGet(target, community, memUsedOID)
	if err != nil {
		return nil, err
	}
	memUsed, _ := strconv.Atoi(memUsedStr)

	// 获取内存空闲信息
	memFreeStr, err := snmp_1_get_client.SnmpGet(target, community, memFreeOID)
	if err != nil {
		return nil, err
	}
	memFree, _ := strconv.Atoi(memFreeStr)

	// 计算内存使用百分比
	memPercent := 0.0
	if memUsed != 0 || memFree != 0 {
		memPercent = float64(memUsed) / float64(memUsed+memFree) * 100
	}

	// 获取接口信息
	ifNameOID := ".1.3.6.1.2.1.2.2.1.2"
	ifStatusOID := ".1.3.6.1.2.1.2.2.1.8"
	ifInBytesOID := ".1.3.6.1.2.1.2.2.1.10"
	ifOutBytesOID := ".1.3.6.1.2.1.2.2.1.16"

	ifNames, err := snmp_2_getbulk_client.SnmpGetBulk(target, community, ifNameOID, 0, 10)
	if err != nil {
		return nil, err
	}

	ifStatuses, err := snmp_2_getbulk_client.SnmpGetBulk(target, community, ifStatusOID, 0, 10)
	if err != nil {
		return nil, err
	}

	ifInBytes, err := snmp_2_getbulk_client.SnmpGetBulk(target, community, ifInBytesOID, 0, 10)
	if err != nil {
		return nil, err
	}

	ifOutBytes, err := snmp_2_getbulk_client.SnmpGetBulk(target, community, ifOutBytesOID, 0, 10)
	if err != nil {
		return nil, err
	}

	interfaces := make([]InterfaceStats, len(ifNames))
	for i := range ifNames {
		status := false
		if ifStatuses[i] == "1" {
			status = true
		}
		inBytes, _ := strconv.Atoi(ifInBytes[i])
		outBytes, _ := strconv.Atoi(ifOutBytes[i])
		interfaces[i] = InterfaceStats{
			Name:     ifNames[i],
			Status:   status,
			InBytes:  inBytes,
			OutBytes: outBytes,
		}
	}

	stats := &DeviceStats{
		CPU:           cpu,
		MemoryUsed:    memUsed,
		MemoryFree:    memFree,
		MemoryPercent: memPercent,
		IPAddress:     target,
		Interfaces:    interfaces,
	}

	return stats, nil
}
