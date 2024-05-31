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
			output += fmt.Sprintf("%s : %s\n", variable.Name, string(variable.Value.([]byte)))
		default:
			output += fmt.Sprintf("%s : %v\n", variable.Name, variable.Value)
		}
	}

	return output, nil
}
