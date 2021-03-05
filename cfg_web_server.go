package config

type WebServer struct {
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	Mode               string `yaml:"mode"`
	HasCORS            bool   `yaml:"has_cors"`
	Compress           bool   `yaml:"compress"`
	Debug              bool   `yaml:"debug"`
	ReadTimeout        int    `yaml:"read_timeout"`
	WriteTimeout       int    `yaml:"write_timeout"`
	IdleTimeout        int    `yaml:"idle_timeout"`
	ShutdownTimeout    int    `yaml:"shutdown_timeout"`
	MaxConnsPerIP      int    `yaml:"max_conn_per_ip"`
	MaxRequestsPerConn int    `yaml:"max_req_per_conn"`
}
