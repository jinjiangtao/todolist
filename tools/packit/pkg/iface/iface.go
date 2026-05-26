package iface

import (
	"github.com/google/gopacket/pcap"
)

type Interface struct {
	Index int
	Name  string
	Desc  string
	IP    string
}

func GetAllInterfaces() ([]Interface, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}

	var ifaces []Interface
	for i, dev := range devices {
		ip := ""
		if len(dev.Addresses) > 0 {
			ip = dev.Addresses[0].IP.String()
		}
		ifaces = append(ifaces, Interface{
			Index: i + 1,
			Name:  dev.Name,
			Desc:  dev.Description,
			IP:    ip,
		})
	}
	return ifaces, nil
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