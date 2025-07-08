package collector

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/http"
)

type TimeSeriesManager struct {
	service http.Service
	client  influxdb2.Client
}

func NewTimeSeriesManager(influxdbUrl, influxdbToken string) *TimeSeriesManager {
	service := http.NewService(influxdbUrl, influxdbToken, http.DefaultOptions())
	client := influxdb2.NewClientWithOptions(influxdbUrl, "c6jGYUCinwzeTWdeUh32", influxdb2.DefaultOptions().SetLogLevel(3))

	s := &TimeSeriesManager{
		service: service,
		client:  client,
	}
	return s
}

func (m *TimeSeriesManager) Init() {

}
