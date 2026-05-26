package packet

import (
	"encoding/binary"
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type PacketInfo struct {
	Timestamp    string
	SrcIP        string
	DstIP        string
	SrcPort      uint16
	DstPort      uint16
	Protocol     string
	Length       int
	ProtocolType string
	Payload      string
}

func ParsePacket(pkt gopacket.Packet) *PacketInfo {
	info := &PacketInfo{}

	if timestamp := pkt.Metadata().Timestamp; !timestamp.IsZero() {
		info.Timestamp = timestamp.Format("15:04:05.000")
	}

	if netLayer := pkt.NetworkLayer(); netLayer != nil {
		info.SrcIP = netLayer.NetworkFlow().Src().String()
		info.DstIP = netLayer.NetworkFlow().Dst().String()
	}

	if tcpLayer := pkt.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		tcp := tcpLayer.(*layers.TCP)
		info.SrcPort = uint16(tcp.SrcPort)
		info.DstPort = uint16(tcp.DstPort)
		info.Protocol = "TCP"
	}

	if udpLayer := pkt.Layer(layers.LayerTypeUDP); udpLayer != nil {
		udp := udpLayer.(*layers.UDP)
		info.SrcPort = uint16(udp.SrcPort)
		info.DstPort = uint16(udp.DstPort)
		info.Protocol = "UDP"
	}

	if icmpLayer := pkt.Layer(layers.LayerTypeICMPv4); icmpLayer != nil {
		info.Protocol = "ICMP"
		info.ProtocolType = "IPv4"
	}
	if icmp6Layer := pkt.Layer(layers.LayerTypeICMPv6); icmp6Layer != nil {
		info.Protocol = "ICMP"
		info.ProtocolType = "IPv6"
	}

	info.Length = len(pkt.Data())

	if dnsLayer := pkt.Layer(layers.LayerTypeDNS); dnsLayer != nil {
		info.Protocol = "DNS"
		dns := dnsLayer.(*layers.DNS)
		if len(dns.Questions) > 0 {
			info.Payload = string(dns.Questions[0].Name)
		}
	}

	if payload := pkt.ApplicationLayer(); payload != nil {
		if info.Protocol == "TCP" {
			data := payload.Payload()
			if len(data) > 0 {
				if data[0] >= 'A' && data[0] <= 'Z' {
					end := 0
					for i, b := range data {
						if b == '\r' || b == '\n' {
							end = i
							break
						}
					}
					if end > 0 {
						info.Payload = string(data[:end])
					} else if len(data) > 64 {
						info.Payload = string(data[:64])
					} else {
						info.Payload = string(data)
					}
				}
			}
		}
	}

	if info.Payload == "" && len(pkt.Data()) > 0 {
		info.Payload = formatPayload(pkt.Data())
	}

	return info
}

func formatPayload(data []byte) string {
	if len(data) > 64 {
		data = data[:64]
	}
	result := ""
	for _, b := range data {
		if b >= 32 && b <= 126 {
			result += string(b)
		} else {
			result += "."
		}
	}
	return result
}

func IsLocalIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	privateBlocks := []net.IPNet{
		{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)},
		{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(12, 32)},
		{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)},
		{IP: net.ParseIP("127.0.0.0"), Mask: net.CIDRMask(8, 32)},
	}
	for _, block := range privateBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

func ParseDNSName(data []byte, offset int) string {
	result := ""
	currentOffset := offset
	for {
		if currentOffset >= len(data) {
			break
		}
		length := int(data[currentOffset])
		if length == 0 {
			currentOffset++
			break
		}
		if (length & 0xC0) == 0xC0 {
			pointerOffset := int(binary.BigEndian.Uint16(data[currentOffset:currentOffset+2]) & 0x3FFF)
			result += ParseDNSName(data, pointerOffset)
			currentOffset += 2
			break
		}
		currentOffset++
		if currentOffset+length > len(data) {
			break
		}
		result += string(data[currentOffset : currentOffset+length]) + "."
		currentOffset += length
	}
	if len(result) > 0 && result[len(result)-1] == '.' {
		result = result[:len(result)-1]
	}
	return result
}