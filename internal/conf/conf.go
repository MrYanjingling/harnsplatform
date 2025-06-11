package conf

import (
	"time"
)

type Bootstrap struct {
	Server *Server `mapstructure:"server,omitempty"`
	Data   *Data   `mapstructure:"data,omitempty"`
}

func (x *Bootstrap) GetServer() *Server {
	if x != nil {
		return x.Server
	}
	return nil
}

func (x *Bootstrap) GetData() *Data {
	if x != nil {
		return x.Data
	}
	return nil
}

type Server struct {
	Http *ServerHTTP `mapstructure:"http,omitempty"`
	Grpc *ServerGRPC `mapstructure:"grpc,omitempty"`
}

func (x *Server) GetHttp() *ServerHTTP {
	if x != nil {
		return x.Http
	}
	return nil
}

func (x *Server) GetGrpc() *ServerGRPC {
	if x != nil {
		return x.Grpc
	}
	return nil
}

type Data struct {
	Database *DataDatabase `mapstructure:"database,omitempty"`
	Redis    *DataRedis    `mapstructure:"redis,omitempty"`
}

func (x *Data) GetDatabase() *DataDatabase {
	if x != nil {
		return x.Database
	}
	return nil
}

func (x *Data) GetRedis() *DataRedis {
	if x != nil {
		return x.Redis
	}
	return nil
}

type ServerHTTP struct {
	Network string        `mapstructure:"network,omitempty"`
	Addr    string        `mapstructure:"addr,omitempty"`
	Timeout time.Duration `mapstructure:"timeout,omitempty"`
}

func (x *ServerHTTP) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *ServerHTTP) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *ServerHTTP) GetTimeout() time.Duration {
	if x != nil {
		return x.Timeout
	}
	return time.Duration(1)
}

type ServerGRPC struct {
	Network string        `mapstructure:"network,omitempty"`
	Addr    string        `mapstructure:"addr,omitempty"`
	Timeout time.Duration `mapstructure:"timeout,omitempty"`
}

func (x *ServerGRPC) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *ServerGRPC) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *ServerGRPC) GetTimeout() time.Duration {
	if x != nil {
		return x.Timeout
	}
	return time.Duration(1)
}

type DataDatabase struct {
	Driver string `mapstructure:"driver,omitempty"`
	Source string `mapstructure:"source,omitempty"`
}

func (x *DataDatabase) GetDriver() string {
	if x != nil {
		return x.Driver
	}
	return ""
}

func (x *DataDatabase) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

type DataRedis struct {
	Network      string        `mapstructure:"network,omitempty"`
	Addr         string        `mapstructure:"addr,omitempty"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout,omitempty"`
	WriteTimeout time.Duration `mapstructure:"write_timeout,omitempty"`
}

func (x *DataRedis) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *DataRedis) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *DataRedis) GetReadTimeout() time.Duration {
	if x != nil {
		return x.ReadTimeout
	}
	return time.Duration(1)
}

func (x *DataRedis) GetWriteTimeout() time.Duration {
	if x != nil {
		return x.WriteTimeout
	}
	return time.Duration(1)
}
