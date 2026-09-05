package internal

type Config struct {
	IpAddr     string
	Port       string
	LeaderAddr string
}

func NewConfig(IpAddr, Port, LeaderAddr string) *Config {
	return &Config{
		IpAddr:     IpAddr,
		Port:       Port,
		LeaderAddr: LeaderAddr,
	}
}

func LoadConfig() *Config {
	return &Config{}
}
