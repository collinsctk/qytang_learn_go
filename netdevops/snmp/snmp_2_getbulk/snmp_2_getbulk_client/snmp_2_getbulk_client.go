package snmp_2_getbulk_client

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"strings"
	"time"
)

// SnmpGetBulk 执行SNMP GETBULK请求并返回结果字符串列表
func SnmpGetBulk(target string, community string, oid string, nonRepeaters uint8, maxRepetitions uint32) ([]string, error) {
	g := gosnmp.Default
	g.Target = target
	g.Port = 161
	g.Community = community
	g.Version = gosnmp.Version2c
	g.Timeout = time.Duration(2) * time.Second

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

	var output []string
	for _, variable := range result.Variables {
		// 过滤掉不是以oid开头的变量
		if strings.HasPrefix(variable.Name, oid) {
			switch variable.Type {
			case gosnmp.OctetString:
				output = append(output, string(variable.Value.([]byte)))
			case gosnmp.Integer, gosnmp.Counter32, gosnmp.Gauge32, gosnmp.TimeTicks, gosnmp.Counter64:
				output = append(output, fmt.Sprintf("%d", gosnmp.ToBigInt(variable.Value).Int64()))
			default:
				output = append(output, fmt.Sprintf("%v", variable.Value))
			}
		}
	}

	return output, nil
}
