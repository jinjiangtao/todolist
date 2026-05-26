package capture

import (
	"fmt"
	"net"
	"strings"
	"sync/atomic"
	"syscall"
	"unsafe"

	"packit/pkg/packet"
)

var (
	ws2_32     = syscall.MustLoadDLL("ws2_32.dll")
	procWSAIoctl = ws2_32.MustFindProc("WSAIoctl")

	SIO_RCVALL = 0x98000001
	RCVALL_ON  uint32 = 1
)

type FilterOptions struct {
	TCP    bool
	UDP    bool
	DNS    bool
	HTTP   bool
	Ports  []uint16
	IP     string
}

type Statistics struct {
	TotalPackets uint64
	TCP          uint64
	UDP          uint64
	DNS          uint64
	HTTP         uint64
	ICMP         uint64
	Other        uint64
}

type Capturer struct {
	socketHandle syscall.Handle
	options      FilterOptions
	stats        Statistics
	running      int32
}

func NewCapturer(ifaceName string, options FilterOptions, outputFile string) (*Capturer, error) {
	sock, err := createRawSocket()
	if err != nil {
		return nil, fmt.Errorf("创建原始套接字失败: %w", err)
	}

	if err := setPromiscuousMode(sock); err != nil {
		syscall.Close(sock)
		return nil, fmt.Errorf("设置混杂模式失败: %w", err)
	}

	if err := bindSocketToInterface(sock, ifaceName); err != nil {
		syscall.Close(sock)
		return nil, fmt.Errorf("绑定网卡失败: %w", err)
	}

	return &Capturer{
		socketHandle: sock,
		options:      options,
		running:      1,
	}, nil
}

func createRawSocket() (syscall.Handle, error) {
	return syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_IP)
}

func setPromiscuousMode(sock syscall.Handle) error {
	var inBuffer uint32 = RCVALL_ON
	var bytesReturned uint32

	ret, _, err := procWSAIoctl.Call(
		uintptr(sock),
		uintptr(SIO_RCVALL),
		uintptr(unsafe.Pointer(&inBuffer)),
		uintptr(4),
		0,
		0,
		uintptr(unsafe.Pointer(&bytesReturned)),
		0,
		0,
	)
	if ret != 0 {
		return err
	}
	return nil
}

func bindSocketToInterface(sock syscall.Handle, ifaceName string) error {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return err
	}

	var ip string
	for _, addr := range addrs {
		ipAddr := addr.String()
		if len(ipAddr) > 0 && ipAddr[0] != ':' {
			idx := strings.Index(ipAddr, "/")
			if idx > 0 {
				ip = ipAddr[:idx]
			} else {
				ip = ipAddr
			}
			break
		}
	}

	if ip == "" {
		return fmt.Errorf("未找到网卡 %s 的IPv4地址", ifaceName)
	}

	ipAddr := net.ParseIP(ip).To4()
	if ipAddr == nil {
		return fmt.Errorf("无效的IPv4地址: %s", ip)
	}

	sockAddr := syscall.SockaddrInet4{Port: 0}
	copy(sockAddr.Addr[:], ipAddr)

	return syscall.Bind(sock, &sockAddr)
}

func (c *Capturer) ApplyBPFFilter() error {
	return nil
}

func (c *Capturer) shouldFilter(info *packet.PacketInfo) bool {
	if c.options.TCP && info.Protocol != "TCP" {
		return true
	}
	if c.options.UDP && info.Protocol != "UDP" {
		return true
	}
	if c.options.DNS && info.Protocol != "DNS" {
		return true
	}
	if c.options.HTTP && info.Protocol != "TCP" {
		return true
	}

	if len(c.options.Ports) > 0 {
		match := false
		for _, p := range c.options.Ports {
			if info.SrcPort == p || info.DstPort == p {
				match = true
				break
			}
		}
		if !match {
			return true
		}
	}

	if c.options.IP != "" {
		if info.SrcIP != c.options.IP && info.DstIP != c.options.IP {
			return true
		}
	}

	return false
}

func (c *Capturer) Capture(callback func(*packet.PacketInfo)) error {
	buf := make([]byte, 65536)

	for atomic.LoadInt32(&c.running) == 1 {
		n, _, err := syscall.Recvfrom(c.socketHandle, buf, 0)
		if err != nil || n == 0 {
			continue
		}

		pkt := packet.ParseRawPacket(buf[:n])
		if pkt.Protocol == "" {
			continue
		}

		if c.shouldFilter(pkt) {
			continue
		}

		atomic.AddUint64(&c.stats.TotalPackets, 1)
		switch pkt.Protocol {
		case "TCP":
			atomic.AddUint64(&c.stats.TCP, 1)
		case "UDP":
			atomic.AddUint64(&c.stats.UDP, 1)
		case "DNS":
			atomic.AddUint64(&c.stats.DNS, 1)
		case "HTTP":
			atomic.AddUint64(&c.stats.HTTP, 1)
		case "ICMP":
			atomic.AddUint64(&c.stats.ICMP, 1)
		default:
			atomic.AddUint64(&c.stats.Other, 1)
		}

		callback(pkt)
	}

	return nil
}

func (c *Capturer) Stop() {
	atomic.StoreInt32(&c.running, 0)
}

func (c *Capturer) Close() {
	syscall.Close(c.socketHandle)
}

func (c *Capturer) GetStats() Statistics {
	return c.stats
}