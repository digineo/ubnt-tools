package discovery

import (
	"fmt"
	"io/ioutil"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func loadFixture(name string) []byte {
	fmt.Printf("Loading fixture 'testdata/%s'\n", name)
	buffer, err := ioutil.ReadFile("testdata/" + name)
	if err != nil {
		panic(err)
	}
	return buffer
}

func packetFromFixture(a *assert.Assertions, name string) *Packet {
	data := loadFixture(name)
	if !a.True(len(data) > 0, "fixture %s is empty", name) {
		return nil
	}

	pkt, err := ParsePacket(data)
	if a.Nil(err) && a.NotNil(pkt) {
		return pkt
	}
	return nil
}

func TestParsePacket(t *testing.T) {
	assert := assert.New(t)

	tt := map[string]string{
		"edgerouter.dat":  "04:18:d6:83:f8:ec",
		"nanobeam-2.dat":  "80:2a:a8:64:a7:12",
		"nanobeam-1b.dat": "80:2a:a8:64:a7:22",
		"nanobeam-1a.dat": "80:2a:a8:64:a7:22",
	}

	for name, mac := range tt {
		dev := packetFromFixture(assert, name)
		for _, t := range dev.Tags {
			if t.ID == tagMacAddress {
				assert.Equal(mac, t.value.(net.HardwareAddr).String())
			}
		}
	}
}

func TestPacketToDevice(t *testing.T) {
	assert := assert.New(t)

	dev := packetFromFixture(assert, "edgerouter.dat").Device()
	const mac = "04:18:d6:83:f8:ec"
	assert.Equal(mac, dev.MacAddress)
	assert.Equal("ERLite-3", dev.Platform)
	assert.Equal("digineo", dev.Hostname)

	if ips, ok := dev.IPAddresses[mac]; assert.True(ok, "expected mac %s to have IP addresses", mac) {
		assert.Contains(ips, "172.16.0.1")
		assert.Contains(ips, "172.16.2.1")
		assert.Contains(ips, "1.2.3.4")
		assert.Contains(ips, "66.66.66.66")
	}
}
