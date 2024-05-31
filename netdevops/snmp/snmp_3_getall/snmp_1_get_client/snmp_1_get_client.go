package snmp_1_get_client

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"time"
)

// SnmpGet 执行SNMP GET请求并返回结果字符串
func SnmpGet(target string, community string, oid string) (string, error) {
	g := gosnmp.Default
	g.Target = target
	g.Port = 161
	g.Community = community
	g.Version = gosnmp.Version2c
	g.Timeout = time.Duration(2) * time.Second

	err := g.Connect()
	if err != nil {
		return "", fmt.Errorf("Connect() err: %v", err)
	}
	defer g.Conn.Close()

	oids := []string{oid}
	result, err := g.Get(oids)
	if err != nil {
		return "", fmt.Errorf("Get() err: %v", err)
	}

	var output string
	for _, variable := range result.Variables {
		switch variable.Type {
		case gosnmp.OctetString:
			output = string(variable.Value.([]byte))
		case gosnmp.Integer, gosnmp.Counter32, gosnmp.Gauge32, gosnmp.TimeTicks, gosnmp.Counter64:
			output = fmt.Sprintf("%d", gosnmp.ToBigInt(variable.Value).Int64())
		default:
			output = fmt.Sprintf("%v", variable.Value)
		}
	}

	return output, nil
}
