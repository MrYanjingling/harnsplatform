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
	Http *Server_HTTP `mapstructure:"http,omitempty"`
	Grpc *Server_GRPC `mapstructure:"grpc,omitempty"`
}

func (x *Server) GetHttp() *Server_HTTP {
	if x != nil {
		return x.Http
	}
	return nil
}

func (x *Server) GetGrpc() *Server_GRPC {
	if x != nil {
		return x.Grpc
	}
	return nil
}

type Data struct {
	Database *Data_Database `mapstructure:"database,omitempty"`
	Redis    *Data_Redis    `mapstructure:"redis,omitempty"`
}

func (x *Data) GetDatabase() *Data_Database {
	if x != nil {
		return x.Database
	}
	return nil
}

func (x *Data) GetRedis() *Data_Redis {
	if x != nil {
		return x.Redis
	}
	return nil
}

type Server_HTTP struct {
	Network string        `mapstructure:"network,omitempty"`
	Addr    string        `mapstructure:"addr,omitempty"`
	Timeout time.Duration `mapstructure:"timeout,omitempty"`
}

func (x *Server_HTTP) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *Server_HTTP) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *Server_HTTP) GetTimeout() time.Duration {
	if x != nil {
		return x.Timeout
	}
	return time.Duration(1)
}

type Server_GRPC struct {
	Network string        `mapstructure:"network,omitempty"`
	Addr    string        `mapstructure:"addr,omitempty"`
	Timeout time.Duration `mapstructure:"timeout,omitempty"`
}

func (x *Server_GRPC) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *Server_GRPC) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *Server_GRPC) GetTimeout() time.Duration {
	if x != nil {
		return x.Timeout
	}
	return time.Duration(1)
}

type Data_Database struct {
	Driver string `mapstructure:"driver,omitempty"`
	Source string `mapstructure:"source,omitempty"`
}

func (x *Data_Database) GetDriver() string {
	if x != nil {
		return x.Driver
	}
	return ""
}

func (x *Data_Database) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

type Data_Redis struct {
	Network      string        `mapstructure:"network,omitempty"`
	Addr         string        `mapstructure:"addr,omitempty"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout,omitempty"`
	WriteTimeout time.Duration `mapstructure:"write_timeout,omitempty"`
}

func (x *Data_Redis) GetNetwork() string {
	if x != nil {
		return x.Network
	}
	return ""
}

func (x *Data_Redis) GetAddr() string {
	if x != nil {
		return x.Addr
	}
	return ""
}

func (x *Data_Redis) GetReadTimeout() time.Duration {
	if x != nil {
		return x.ReadTimeout
	}
	return time.Duration(1)
}

func (x *Data_Redis) GetWriteTimeout() time.Duration {
	if x != nil {
		return x.WriteTimeout
	}
	return time.Duration(1)
}
