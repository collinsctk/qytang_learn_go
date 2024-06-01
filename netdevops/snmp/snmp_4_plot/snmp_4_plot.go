package main

import (
	"encoding/json"
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"log"
	"snmp_4_plot/plot_cpu"
	"snmp_4_plot/smtp_sendmail"
)

func main() {
	// 查询数据并绘制图表
	imagePath, err := queryAndPlotData()
	if err != nil {
		log.Fatalf("Error querying and plotting data: %v", err)
	}

	// 发送带有附件的电子邮件
	err = smtp_sendmail.SendEmailWithAttachment(
		"smtp.qq.com", "465", "3348326959@qq.com", "anchwprpwxfbdbif",
		"3348326959@qq.com", "collinsctk@qytang.com",
		"CPU Usage Plot", "Please find the attached CPU usage plot.",
		imagePath,
	)
	if err != nil {
		log.Fatalf("Error sending email: %v", err)
	}
}

func queryAndPlotData() (string, error) {
	// 连接到 InfluxDB
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://10.10.1.200:8086",
		Username: "qytdbuser",
		Password: "Cisc0123",
	})
	if err != nil {
		return "", fmt.Errorf("error creating InfluxDB client: %v", err)
	}
	defer c.Close()

	// 查询 InfluxDB 数据
	q := client.NewQuery("SELECT mean(\"cpu_usage\") FROM \"router_monitor\" WHERE time > now() - 5m GROUP BY time(1m), \"device_ip\" fill(0)", "qytdb", "s")
	response, err := c.Query(q)
	if err != nil || response.Error() != nil {
		return "", fmt.Errorf("error querying InfluxDB: %v", err)
	}

	// 解析查询结果
	deviceData := make(map[string]plot_cpu.XYs)
	for _, result := range response.Results {
		for _, series := range result.Series {
			deviceIP := series.Tags["device_ip"]
			for _, row := range series.Values {
				timestampFloat, ok := row[0].(json.Number)
				if !ok {
					return "", fmt.Errorf("error parsing timestamp: %v", row[0])
				}
				timestamp, err := timestampFloat.Float64()
				if err != nil {
					return "", fmt.Errorf("error converting timestamp to float64: %v", err)
				}

				valueStr, ok := row[1].(json.Number)
				if !ok {
					return "", fmt.Errorf("error parsing value: %v", row[1])
				}
				valueFloat, err := valueStr.Float64()
				if err != nil {
					return "", fmt.Errorf("error converting value to float64: %v", err)
				}
				deviceData[deviceIP] = append(deviceData[deviceIP], plot_cpu.XY{X: timestamp, Y: valueFloat})
			}
		}
	}

	// 调用绘图函数
	imagePath := "/qyt_learn_go/netdevops/snmp/snmp_4_plot/cpu_usage.png"
	err = plot_cpu.PlotData(deviceData, "CPU Usage Over Time", "Time", "CPU Usage (%)", imagePath)
	if err != nil {
		return "", fmt.Errorf("error plotting data: %v", err)
	}

	fmt.Println("Plot saved to", imagePath)
	return imagePath, nil
}
