package capture

import (
	"os"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"packit/pkg/packet"
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
	handle      *pcap.Handle
	pcapWriter  *pcapgo.Writer
	options     FilterOptions
	stats       Statistics
	running     int32
	outputFile  *os.File
}

func NewCapturer(ifaceName string, options FilterOptions, outputFile string) (*Capturer, error) {
	handle, err := pcap.OpenLive(ifaceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	var writer *pcapgo.Writer
	var f *os.File
	if outputFile != "" {
		f, err = os.Create(outputFile)
		if err != nil {
			handle.Close()
			return nil, err
		}
		writer = pcapgo.NewWriter(f)
		writer.WriteFileHeader(65536, handle.LinkType())
	}

	return &Capturer{
		handle:      handle,
		pcapWriter:  writer,
		options:     options,
		running:     1,
		outputFile:  f,
	}, nil
}

func (c *Capturer) ApplyBPFFilter() error {
	filterParts := []string{}

	if c.options.TCP {
		filterParts = append(filterParts, "tcp")
	}
	if c.options.UDP {
		filterParts = append(filterParts, "udp")
	}
	if c.options.DNS {
		filterParts = append(filterParts, "udp port 53")
	}
	if c.options.HTTP {
		filterParts = append(filterParts, "(tcp port 80 or tcp port 8080)")
	}

	if len(c.options.Ports) > 0 {
		portStrs := []string{}
		for _, p := range c.options.Ports {
			portStrs = append(portStrs, strconv.Itoa(int(p)))
		}
		filterParts = append(filterParts, "port ("+strings.Join(portStrs, " or ")+")")
	}

	if c.options.IP != "" {
		filterParts = append(filterParts, "host "+c.options.IP)
	}

	if len(filterParts) > 0 {
		filter := strings.Join(filterParts, " and ")
		if err := c.handle.SetBPFFilter(filter); err != nil {
			return err
		}
	}

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
	packetSource := gopacket.NewPacketSource(c.handle, c.handle.LinkType())
	for pkt := range packetSource.Packets() {
		if atomic.LoadInt32(&c.running) == 0 {
			break
		}

		info := packet.ParsePacket(pkt)
		if info.Protocol == "" {
			continue
		}

		if c.shouldFilter(info) {
			continue
		}

		atomic.AddUint64(&c.stats.TotalPackets, 1)
		switch info.Protocol {
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

		if c.pcapWriter != nil {
			c.pcapWriter.WritePacket(pkt.Metadata().CaptureInfo, pkt.Data())
		}

		callback(info)
	}
	return nil
}

func (c *Capturer) Stop() {
	atomic.StoreInt32(&c.running, 0)
}

func (c *Capturer) Close() {
	c.handle.Close()
	if c.outputFile != nil {
		c.outputFile.Close()
	}
}

func (c *Capturer) GetStats() Statistics {
	return c.stats
}