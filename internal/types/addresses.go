package types

import (
	"fmt"
	"net"
)

// Address сетевой адрес (хост, порт)
type Address struct {
	Host string
	Port int
}

// String строковое представление
func (a Address) String() string {
	if ip := net.ParseIP(a.Host); ip != nil && ip.To4() == nil {
		return fmt.Sprintf("[%s]:%d", a.Host, a.Port)
	}
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
