package packet

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	"time"
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

type ipHeader struct {
	Version  uint8
	IHL      uint8
	TOS      uint8
	TotalLen uint16
	ID       uint16
	Flags    uint16
	TTL      uint8
	Protocol uint8
	Checksum uint16
	SrcIP    [4]byte
	DstIP    [4]byte
}

type tcpHeader struct {
	SrcPort   uint16
	DstPort   uint16
	Seq       uint32
	Ack       uint32
	DataOffset uint8
	Reserved  uint8
	Flags     uint16
	Window    uint16
	Checksum  uint16
	UrgPtr    uint16
}

type udpHeader struct {
	SrcPort  uint16
	DstPort  uint16
	Length   uint16
	Checksum uint16
}

type icmpHeader struct {
	Type     uint8
	Code     uint8
	Checksum uint16
}

func ParseRawPacket(data []byte) *PacketInfo {
	info := &PacketInfo{}

	info.Timestamp = time.Now().Format("15:04:05.000")

	if len(data) < 20 {
		return info
	}

	ipHdr := parseIPHeader(data)
	info.SrcIP = net.IP(ipHdr.SrcIP[:]).String()
	info.DstIP = net.IP(ipHdr.DstIP[:]).String()
	info.Length = int(ipHdr.TotalLen)

	ipDataOffset := int(ipHdr.IHL * 4)
	transportData := data[ipDataOffset:]

	switch ipHdr.Protocol {
	case 6:
		info.Protocol = "TCP"
		parseTCP(transportData, info)
	case 17:
		info.Protocol = "UDP"
		parseUDP(transportData, info)
	case 1:
		info.Protocol = "ICMP"
		parseICMP(transportData, info)
	default:
		info.Protocol = "OTHER"
	}

	if info.Payload == "" && len(transportData) > 0 {
		info.Payload = formatPayload(transportData)
	}

	return info
}

func parseIPHeader(data []byte) ipHeader {
	hdr := ipHeader{}
	hdr.Version = (data[0] >> 4) & 0x0F
	hdr.IHL = data[0] & 0x0F
	hdr.TOS = data[1]
	hdr.TotalLen = binary.BigEndian.Uint16(data[2:4])
	hdr.ID = binary.BigEndian.Uint16(data[4:6])
	hdr.Flags = binary.BigEndian.Uint16(data[6:8])
	hdr.TTL = data[8]
	hdr.Protocol = data[9]
	hdr.Checksum = binary.BigEndian.Uint16(data[10:12])
	copy(hdr.SrcIP[:], data[12:16])
	copy(hdr.DstIP[:], data[16:20])
	return hdr
}

func parseTCP(data []byte, info *PacketInfo) {
	if len(data) < 20 {
		return
	}

	tcpHdr := tcpHeader{}
	tcpHdr.SrcPort = binary.BigEndian.Uint16(data[0:2])
	tcpHdr.DstPort = binary.BigEndian.Uint16(data[2:4])
	tcpHdr.Seq = binary.BigEndian.Uint32(data[4:8])
	tcpHdr.Ack = binary.BigEndian.Uint32(data[8:12])
	tcpHdr.DataOffset = (data[12] >> 4) & 0x0F

	info.SrcPort = tcpHdr.SrcPort
	info.DstPort = tcpHdr.DstPort

	dataOffset := int(tcpHdr.DataOffset * 4)
	if len(data) > dataOffset {
		payload := data[dataOffset:]
		info.Payload = formatPayload(payload)

		if tcpHdr.DstPort == 80 || tcpHdr.DstPort == 8080 || tcpHdr.SrcPort == 80 || tcpHdr.SrcPort == 8080 {
			if len(payload) > 0 && payload[0] >= 'A' && payload[0] <= 'Z' {
				end := strings.Index(string(payload), "\r\n")
				if end > 0 {
					info.Payload = string(payload[:end])
				}
			}
		}
	}
}

func parseUDP(data []byte, info *PacketInfo) {
	if len(data) < 8 {
		return
	}

	udpHdr := udpHeader{}
	udpHdr.SrcPort = binary.BigEndian.Uint16(data[0:2])
	udpHdr.DstPort = binary.BigEndian.Uint16(data[2:4])
	udpHdr.Length = binary.BigEndian.Uint16(data[4:6])

	info.SrcPort = udpHdr.SrcPort
	info.DstPort = udpHdr.DstPort

	if udpHdr.SrcPort == 53 || udpHdr.DstPort == 53 {
		info.Protocol = "DNS"
		if len(data) > 8 {
			dnsData := data[8:]
			if len(dnsData) > 12 {
				qname := parseDNSName(dnsData, 12)
				if qname != "" {
					info.Payload = qname
				}
			}
		}
	} else if len(data) > 8 {
		payload := data[8:]
		info.Payload = formatPayload(payload)
	}
}

func parseICMP(data []byte, info *PacketInfo) {
	if len(data) < 4 {
		return
	}

	icmpHdr := icmpHeader{}
	icmpHdr.Type = data[0]
	icmpHdr.Code = data[1]

	switch icmpHdr.Type {
	case 0:
		info.Payload = "Echo Reply"
	case 8:
		info.Payload = "Echo Request"
	case 3:
		info.Payload = "Destination Unreachable"
	case 11:
		info.Payload = "Time Exceeded"
	default:
		info.Payload = fmt.Sprintf("Type %d Code %d", icmpHdr.Type, icmpHdr.Code)
	}
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

func parseDNSName(data []byte, offset int) string {
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
			result += parseDNSName(data, pointerOffset)
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