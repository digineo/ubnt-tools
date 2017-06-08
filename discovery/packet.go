package discovery

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

// Packet holds information about a UBNT Discovery response
type Packet struct {
	Version   uint8
	Tags      []*Tag
	timestamp time.Time
}

// ParsePacket tries to parse UPD packet data into a Packet
func ParsePacket(raw []byte) (*Packet, error) {
	if len(raw) <= 4 {
		return nil, fmt.Errorf("packet data too short (%d bytes)", len(raw))
	}

	ver := uint8(raw[0])
	cmd := uint8(raw[1])
	length := binary.BigEndian.Uint16(raw[2:4])

	if int(length)+4 != len(raw) {
		return nil, fmt.Errorf("packet length mismatch (expected %d bytes, got %d)", length+4, len(raw))
	}

	p := &Packet{
		Version:   ver,
		timestamp: time.Now(),
	}
	if err := p.parse(cmd, raw[4:length+4]); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Packet) parse(cmd uint8, data []byte) error {
	if !(p.Version == 1 && cmd == 0) && p.Version != 2 {
		return fmt.Errorf("unsupported packet ver=%d cmd=%d", p.Version, cmd)
	}

	for curr := 0; curr < len(data); {
		id := TagID(data[curr+0])
		n := binary.BigEndian.Uint16(data[curr+1 : curr+3])
		begin, end := curr+3, curr+3+int(n)

		tag, err := ParseTag(id, n, data[begin:end])
		if err != nil {
			log.Print(err)
		} else {
			p.Tags = append(p.Tags, tag)
		}

		curr = end
	}

	return nil
}

// Device converts the packet information into a new device
func (p *Packet) Device() *Device {
	dev := &Device{
		IPAddresses: make(map[string][]string),
		LastSeenAt:  p.timestamp,
		FirstSeenAt: time.Now(),
	}

	for _, t := range p.Tags {
		switch t.ID {
		case tagModelV1, tagModelV2:
			t.StringInto(&dev.Model)
		case tagPlatform:
			t.StringInto(&dev.Platform)
		case tagFirmware:
			t.StringInto(&dev.Firmware)
		case tagEssid:
			t.StringInto(&dev.Essid)
		case tagHostname:
			t.StringInto(&dev.Hostname)

		case tagMacAddress:
			if v, ok := t.value.(net.HardwareAddr); ok {
				dev.MacAddress = v.String()
			}
		case tagUptime:
			if v, ok := t.value.(uint32); ok {
				now := time.Now()
				dur := -1 * int(v)
				dev.UpSince = now.Add(time.Duration(dur) * time.Second)
			}
		case tagWmode:
			if v, ok := t.value.(uint8); ok {
				switch v {
				case 2:
					dev.WirelessMode = "Station"
				case 3:
					dev.WirelessMode = "AccessPoint"
				default:
					dev.WirelessMode = fmt.Sprintf("unknown (%#02x)", v)
				}
			}
		case tagIPInfo:
			if v, ok := t.value.(*ipInfo); ok {
				m := v.MacAddress.String()
				dev.IPAddresses[m] = append(dev.IPAddresses[m], v.IPAddress.String())
			}
		}
	}
	return dev
}
