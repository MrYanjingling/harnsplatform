package conf

type BrokerBootstrap struct {
	Config *BrokerConfig   `mapstructure:"config,omitempty"`
	Server *Server         `mapstructure:"server,omitempty"`
	Data   *TimeSeriesData `mapstructure:"data,omitempty"`
}

type TimeSeriesData struct {
	Influxdb *InfluxDb  `mapstructure:"influxdb,omitempty"`
	Redis    *DataRedis `mapstructure:"redis,omitempty"`
}

type InfluxDb struct {
	Url   string `mapstructure:"url,omitempty"`
	Token string `mapstructure:"token,omitempty"`
}

func (x *InfluxDb) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *InfluxDb) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type BrokerConfig struct {
	TimeSeriesStore TimeSeriesStorePeriod `mapstructure:"timeSeriesStore,omitempty"`
	Sink            Sink                  `mapstructure:"sink,omitempty"`
}

type TimeSeriesStorePeriod struct {
	Flag     bool   `mapstructure:"flag,omitempty"`
	TimeType string `mapstructure:"timeType,omitempty"`
	Period   int    `mapstructure:"period,omitempty"`
}

func (x *TimeSeriesStorePeriod) GetFlag() bool {
	if x != nil {
		return x.Flag
	}
	return false
}

func (x *TimeSeriesStorePeriod) GetTimeType() string {
	if x != nil {
		return x.TimeType
	}
	return ""
}

func (x *TimeSeriesStorePeriod) GetPeriod() int {
	if x != nil {
		return x.Period
	}
	return 0
}

type Sink struct {
	Flag          bool                   `mapstructure:"flag,omitempty"`
	SinkMQ        string                 `mapstructure:"sinkMQ,omitempty"` // kafka  mqtt
	MQConfig      map[string]interface{} `mapstructure:"mqConfig,omitempty"`
	PushFrequency int                    `mapstructure:"pushFrequency,omitempty"` // 单位 秒
}

func (x *Sink) GetFlag() bool {
	if x != nil {
		return x.Flag
	}
	return false
}

func (x *Sink) GetSinkMQ() string {
	if x != nil {
		return x.SinkMQ
	}
	return ""
}
func (x *Sink) GetMQConfig() map[string]interface{} {
	if x != nil {
		return x.MQConfig
	}
	return nil
}
func (x *Sink) GetPushFrequency() int {
	if x != nil {
		return x.PushFrequency
	}
	return 0
}
