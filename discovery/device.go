package discovery

import "time"

// Device descibes an UBNT device found on the local network
type Device struct {
	Model        string
	Platform     string
	MacAddress   string
	Hostname     string
	Firmware     string
	IPAddresses  map[string][]string
	UpSince      time.Time
	Essid        string
	WirelessMode string
	LastSeenAt   time.Time
	FirstSeenAt  time.Time
}

// RecentlySeen tells you, whether you have seen this device in the
// given time period
func (d *Device) RecentlySeen(dur time.Duration) bool {
	return d.LastSeenAt.Add(dur).After(time.Now())
}

// Merge updates this instance with the values of the other (by copying
// the data), so that references to this instance are kept intact
func (d *Device) Merge(other *Device) {
	d.Model = other.Model
	d.Platform = other.Platform
	d.MacAddress = other.MacAddress
	d.Hostname = other.Hostname
	d.Firmware = other.Firmware
	d.IPAddresses = make(map[string][]string)
	for mac, ips := range other.IPAddresses {
		d.IPAddresses[mac] = append(d.IPAddresses[mac], ips...)
	}
	d.UpSince = other.UpSince
	d.Essid = other.Essid
	d.WirelessMode = other.WirelessMode
	d.LastSeenAt = other.LastSeenAt
	if d.FirstSeenAt.After(other.FirstSeenAt) {
		d.FirstSeenAt = other.FirstSeenAt
	}
}

// Clone creates a deep-copy
func (d *Device) Clone() *Device {
	other := &Device{FirstSeenAt: time.Now()}
	other.Merge(d)
	return other
}

func (d *Device) String() string {
	buf := "Device"
	buf += "\n  MAC:          " + d.MacAddress
	buf += "\n  Model:        " + d.Model
	buf += "\n  Platform:     " + d.Platform
	buf += "\n  Firmware:     " + d.Firmware
	buf += "\n  Hostname:     " + d.Hostname

	now := time.Now()
	if (d.UpSince != time.Time{}) {
		buf += "\n  Booted at:    " + d.UpSince.Format(time.RFC3339)
		buf += "\n  booted:       " + now.Sub(d.UpSince).String() + " ago"
	}

	buf += "\n  first seen:   " + now.Sub(d.FirstSeenAt).String() + " ago"
	buf += "\n  last seen:    " + now.Sub(d.FirstSeenAt).String() + " ago"

	for mac, ips := range d.IPAddresses {
		buf += "\n  IP addresses on interface " + mac
		for _, ip := range ips {
			buf += "\n    - " + ip
		}
	}

	if d.Essid != "" {
		buf += "\n  ESSID:        " + d.Essid
	}
	if d.WirelessMode != "" {
		buf += "\n  WMode:        " + d.WirelessMode
	}

	return buf
}
