package discovery

import (
	"encoding/binary"
	"fmt"
	"net"
)

// TagID identifies tags
type TagID uint8

const (
	// common tags
	tagMacAddress = 0x01 // mac addr
	tagIPInfo     = 0x02 // mac addr + ipv4 addr
	tagFirmware   = 0x03 // string
	tagUptime     = 0x0A // uint32
	tagHostname   = 0x0B // string
	tagPlatform   = 0x0C // string
	tagEssid      = 0x0D // string
	tagWmode      = 0x0E // uint8

	// v1 tags
	tagUsername     = 0x06 // string
	tagSalt         = 0x07 // ?
	tagRndChallenge = 0x08 // ?
	tagChallenge    = 0x09 // ?
	tagWebui        = 0x0F // string?
	tagModelV1      = 0x14 // string

	// v2 tags
	tagSequence     = 0x12 // uint?
	tagSourceMac    = 0x13 // string
	tagShortVersion = 0x16 // string
	tagDefault      = 0x17 // uint8 (bool)
	tagLocating     = 0x18 // uint8 (bool)
	tagDhcpc        = 0x19 // uint8 (bool)
	tagDhcpcBound   = 0x1A // uint8 (bool)
	tagReqFirmware  = 0x1B // string
	tagSshdPort     = 0x1C // uint32?
	tagModelV2      = 0x15 // string
)

type tagParser func([]byte) (interface{}, error)

// TagDescription annotates some meta information to a TagID
type TagDescription struct {
	shortName string
	longName  string
	byteLen   int // 0 <= unspec, 0 < length < 2**16, 2**16 >= error
	converter tagParser
}

var (
	tagDescriptions = map[TagID]TagDescription{
		tagEssid:        {"essid", "Wireless ESSID", -1, parseString},
		tagFirmware:     {"firmware", "Firmware", -1, parseString},
		tagHostname:     {"hostname", "Hostname", -1, parseString},
		tagIPInfo:       {"ipinfo", "MAC/IP mapping", 10, parseIPInfo},
		tagMacAddress:   {"hwaddr", "Hardware/MAC address", 6, parseMacAddress},
		tagModelV1:      {"model.v1", "Model name", -1, parseString},
		tagModelV2:      {"model.v2", "Model name", -1, parseString},
		tagPlatform:     {"platform", "Platform information", -1, parseString},
		tagShortVersion: {"short-ver", "Short version", -1, parseString},
		tagSshdPort:     {"sshd-port", "SSH port", 2, parseUint16},
		tagUptime:       {"uptime", "Uptime", 4, parseUint32},
		tagUsername:     {"username", "Username", -1, parseString},
		tagWebui:        {"webui", "URL for Web-UI", -1, nil},
		tagWmode:        {"wmode", "Wireless mode", 1, parseUint8},

		// unknown or not yet found in the wild
		tagChallenge:    {"challenge", "(?)", -1, nil},
		tagDefault:      {"default", "(bool)", 1, parseBool},
		tagDhcpc:        {"dhcpc", "(bool)", 1, parseBool},
		tagDhcpcBound:   {"dhcpc-bound", "(bool)", 1, parseBool},
		tagLocating:     {"locating", "(bool)", 1, parseBool},
		tagReqFirmware:  {"req-firmware", "(string)", -1, parseString},
		tagRndChallenge: {"rnd-challenge", "(?)", -1, nil},
		tagSalt:         {"salt", "(?)", -1, nil},
		tagSequence:     {"seq", "(uint?)", -1, nil},
		tagSourceMac:    {"source-mac", "(?)", -1, nil},
	}
)

// Tag describes a key value pair
type Tag struct {
	ID          TagID
	description *TagDescription
	value       interface{}
}

type ipInfo struct {
	MacAddress net.HardwareAddr
	IPAddress  net.IP
}

// ParseTag converts a byte stream (i.e. an UDP packet slice) into a Tag
func ParseTag(id TagID, n uint16, raw []byte) (*Tag, error) {
	t := &Tag{ID: id}

	// check if known, unknown, or not yet seen
	if d, ok := tagDescriptions[t.ID]; ok {
		t.description = &d
	} else {
		t.description = &TagDescription{
			shortName: "unknown",
			longName:  fmt.Sprintf("unknown (%#x)", t.ID),
		}
	}

	if t.description.byteLen > 0 {
		if t.description.byteLen != int(n) {
			return nil, fmt.Errorf(
				"length mismatch for tag %s (expected %d bytes, got %d)",
				t.description.shortName, t.description.byteLen, n,
			)
		}
		if len(raw) < int(n) {
			return nil, fmt.Errorf(
				"not enough data for tag %s (expected %d bytes, got %d)",
				t.description.shortName, n, len(raw),
			)
		}
	}

	if val, err := t.description.convert(raw); err == nil {
		t.value = val
	} else {
		return nil, err
	}

	return t, nil
}

// Name returns the short tag name
func (t *Tag) Name() string {
	return t.description.shortName
}

// Description returns the long tag name
func (t *Tag) Description() string {
	return t.description.longName
}

// StringInto tries to update the given string reference with a type
// asserted value (it doesn't perform an update, if the type assertion
// fails)
func (t *Tag) StringInto(ref *string) {
	if v, ok := t.value.(string); ok {
		*ref = v
	}
}

func (td *TagDescription) convert(data []byte) (interface{}, error) {
	if td.converter == nil {
		return fmt.Sprintf("len:%d<%x>", len(data), data), nil
	}
	return td.converter(data)
}

func parseString(data []byte) (interface{}, error) {
	return string(data), nil
}

func parseBool(data []byte) (interface{}, error) {
	return uint8(data[0]) != 0, nil
}

func parseUint8(data []byte) (interface{}, error) {
	return uint8(data[0]), nil
}

func parseUint16(data []byte) (interface{}, error) {
	return binary.BigEndian.Uint16(data[0:2]), nil
}

func parseUint32(data []byte) (interface{}, error) {
	return binary.BigEndian.Uint32(data[0:4]), nil
}

func parseMacAddress(data []byte) (interface{}, error) {
	return net.HardwareAddr(data[0:6]), nil
}

func parseIPInfo(data []byte) (interface{}, error) {
	return &ipInfo{
		MacAddress: net.HardwareAddr(data[0:6]),
		IPAddress:  net.IP(data[6:10]),
	}, nil
}
