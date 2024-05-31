package influxdb_client

import (
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

type InfluxDBClient struct {
	client   client.Client
	database string
}

func NewInfluxDBClient(host, port, username, password, database string) (*InfluxDBClient, error) {
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://" + host + ":" + port,
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, err
	}

	return &InfluxDBClient{client: cli, database: database}, nil
}

func (i *InfluxDBClient) WritePoints(measurement string, tags map[string]string, fields map[string]interface{}) error {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  i.database,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	pt, err := client.NewPoint(measurement, tags, fields, time.Now())
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	return i.client.Write(bp)
}

func (i *InfluxDBClient) Close() {
	i.client.Close()
}
