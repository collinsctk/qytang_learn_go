package snmp_2_getbulk_client

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"time"
)

// SnmpGetBulk 执行SNMP GETBULK请求并返回结果字符串数组
func SnmpGetBulk(target string, community string, oid string, nonRepeaters uint8, maxRepetitions uint32) ([]string, error) {
	g := &gosnmp.GoSNMP{
		Target:    target,
		Port:      161,
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		Retries:   1,
	}

	err := g.Connect()
	if err != nil {
		return nil, fmt.Errorf("Connect() err: %v", err)
	}
	defer g.Conn.Close()

	oids := []string{oid}
	result, err := g.GetBulk(oids, nonRepeaters, maxRepetitions)
	if err != nil {
		return nil, fmt.Errorf("GetBulk() err: %v", err)
	}

	var outputs []string
	for _, variable := range result.Variables {
		switch variable.Type {
		case gosnmp.OctetString:
			outputs = append(outputs, string(variable.Value.([]byte)))
		case gosnmp.Integer, gosnmp.Counter32, gosnmp.Gauge32, gosnmp.TimeTicks, gosnmp.Counter64:
			outputs = append(outputs, fmt.Sprintf("%d", gosnmp.ToBigInt(variable.Value).Int64()))
		default:
			outputs = append(outputs, fmt.Sprintf("%v", variable.Value))
		}
	}

	return outputs, nil
}
