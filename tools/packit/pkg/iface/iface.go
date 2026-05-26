package iface

import (
	"fmt"
	"net"
	"strings"
)

type Interface struct {
	Index int
	Name  string
	Desc  string
	IP    string
}

func GetAllInterfaces() ([]Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("获取网络接口失败: %w", err)
	}

	var result []Interface
	for i, iface := range ifaces {
		ip := ""
		addrs, err := iface.Addrs()
		if err == nil && len(addrs) > 0 {
			for _, addr := range addrs {
				ipAddr := addr.String()
				if len(ipAddr) > 0 {
					idx := strings.Index(ipAddr, "/")
					if idx > 0 {
						ip = ipAddr[:idx]
					} else {
						ip = ipAddr
					}
					if strings.Contains(ip, ".") {
						break
					}
				}
			}
		}

		desc := iface.Name
		if iface.Flags&net.FlagUp != 0 {
			desc += " (UP)"
		}

		result = append(result, Interface{
			Index: i + 1,
			Name:  iface.Name,
			Desc:  desc,
			IP:    ip,
		})
	}
	return result, nil
}

func GetInterfaceByIndex(index int) (*Interface, error) {
	ifaces, err := GetAllInterfaces()
	if err != nil {
		return nil, err
	}
	if index < 1 || index > len(ifaces) {
		return nil, nil
	}
	return &ifaces[index-1], nil
}

func GetInterfaceByName(name string) (*Interface, error) {
	ifaces, err := GetAllInterfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Name == name {
			return &iface, nil
		}
	}
	return nil, nil
}