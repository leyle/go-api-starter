package confighelper

import "fmt"

// generic config model
type ConnectionOption struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Passwd   string `json:"passwd" yaml:"passwd"`
	Protocol string `json:"protocol" yaml:"protocol"`
}

func (c *ConnectionOption) ListenServerAddr() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
