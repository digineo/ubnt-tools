package web

import "github.com/digineo/ubnt-tools/provisioner"

// DeviceJSON wraps a provisioner.Device into JSON presentation
type DeviceJSON struct {
	CanUpgrade   bool                `json:"can_upgrade"`
	Essid        string              `json:"essid"`
	Firmware     string              `json:"firmware"`
	FirstSeenAt  int64               `json:"first_seen_at"`
	HasConfig    bool                `json:"has_config"`
	Hostname     string              `json:"hostname"`
	IPAddress    string              `json:"ip_address"`
	IPAddresses  map[string][]string `json:"ip_addresses"`
	LastSeenAt   int64               `json:"last_seen_at"`
	MacAddress   string              `json:"mac_address"`
	Model        string              `json:"model"`
	Platform     string              `json:"platform"`
	Status       string              `json:"status"`
	UpSince      int64               `json:"up_since"`
	WirelessMode string              `json:"wireless_mode"`
}

// MakeDeviceJSON transforms a Device into a DeviceJSON
func MakeDeviceJSON(dev *provisioner.Device) *DeviceJSON {
	j := &DeviceJSON{
		CanUpgrade:   dev.CanUpgrade(),
		Essid:        dev.Essid,
		Firmware:     dev.Firmware,
		FirstSeenAt:  dev.FirstSeenAt.Unix(),
		HasConfig:    dev.HasConfig(),
		Hostname:     dev.Hostname,
		IPAddress:    dev.IPAddress,
		IPAddresses:  make(map[string][]string),
		LastSeenAt:   dev.LastSeenAt.Unix(),
		MacAddress:   dev.MacAddress,
		Model:        dev.Model,
		Platform:     dev.Platform,
		Status:       dev.Status(),
		UpSince:      dev.UpSince.Unix(),
		WirelessMode: dev.WirelessMode,
	}

	for mac, ips := range dev.IPAddresses {
		// copy(j.IPAddresses[mac], ips) // doesn't work

		j.IPAddresses[mac] = make([]string, len(ips), len(ips))
		for i, ip := range ips {
			j.IPAddresses[mac][i] = ip
		}
	}

	return j
}

// WrapDeviceJSON transforms a list of Devices into a list of DeviceJSONs
func WrapDeviceJSON(devices []*provisioner.Device) []*DeviceJSON {
	list := make([]*DeviceJSON, len(devices))
	for i, d := range devices {
		list[i] = MakeDeviceJSON(d)
	}
	return list
}
