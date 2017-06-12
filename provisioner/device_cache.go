package provisioner

import (
	"log"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/digineo/goldflags"
	"github.com/digineo/ubnt-tools/discovery"
)

type deviceCache struct {
	list    map[string]*Device
	updated time.Time
	sync.RWMutex
}

// StartAutoDiscover starts the UBNT auto discovery mechanism. See
// discovery.AutoDiscover for details.
func (c *Configuration) StartAutoDiscover(notify discovery.NotifyHandler) (d *discovery.Discover, err error) {
	d, err = discovery.AutoDiscover(notify, c.InterfaceNames...)
	if err == nil {
		c.autoDiscoverer = d
	}
	return
}

// GetDevices returns an array with all discovered devices.
func (c *Configuration) GetDevices() (list []*Device) {
	c.updateCache()
	c.devices.RLock()
	for _, dev := range c.devices.list {
		list = append(list, dev)
	}
	c.devices.RUnlock()
	return
}

// FindDevice searches the list of discovered devices and returns a
// pointer to it (or nil, if we can't find it).
func (c *Configuration) FindDevice(mac string) *Device {
	c.updateCache()
	c.devices.RLock()
	defer c.devices.RUnlock()

	if dev, found := c.devices.list[mac]; found {
		return dev
	}
	return nil
}

func (c *Configuration) updateCache() {
	c.devices.Lock()
	defer c.devices.Unlock()

	if time.Since(c.devices.updated) < 1*time.Second {
		return
	}
	if c.autoDiscoverer == nil {
		log.Printf("[device cache] no auto-discoverer found")
		return
	}

	log.Printf("[device cache] updating")
	discovered := c.autoDiscoverer.List()
	seen := make(map[string]int) // IP address -> # of devices with this address
	list := c.devices.list       // mac -> Device

	for _, dev := range discovered {
		for _, addrs := range dev.IPAddresses {
			for _, ip := range addrs {
				seen[ip]++
			}
		}
		if old, found := list[dev.MacAddress]; found {
			old.Device.Merge(dev)
			old.firmwarePath = ""
			old.systemConfigPath = ""
		} else {
			list[dev.MacAddress] = &Device{Device: dev}
		}
	}

	// inject additional information
	for _, dev := range list {
		// unique IP addresses
		for _, addrs := range dev.IPAddresses {
			dev.IPAddress = ""
			for _, ip := range addrs {
				if ip == "192.168.1.20" {
					continue
				}
				if dev.IPAddress == "" && seen[ip] == 1 {
					dev.IPAddress = ip
				}
			}
		}

		// Path to system config
		if cfgPath := filepath.Join(c.ConfigDirectory, sanitizeMac(dev.MacAddress)+".cfg"); goldflags.PathExist(cfgPath) {
			dev.systemConfigPath = cfgPath
		}

		// Path to firmware
		if target, ok := c.reverseUpgradePaths[dev.Firmware]; ok {
			if fwPath := filepath.Join(c.FirmwareDirectory, target); goldflags.PathExist(fwPath) {
				dev.firmwarePath = fwPath
			}
		}

		// SSH auth methods
		dev.authMethods = c.sshAuthMethods
	}

	c.devices.updated = time.Now()
	return
}

func sanitizeMac(mac string) string {
	return strings.ToLower(strings.Replace(mac, ":", "", -1))
}
