package web

import (
	"net/http"

	"github.com/digineo/ubnt-tools/provisioner"
	"github.com/gorilla/mux"
)

// GET /api/devices
func (g *goWeb) getDevices(w http.ResponseWriter, r *http.Request) {
	devices := WrapDeviceJSON(g.config.GetDevices())
	g.responseJSON(w, http.StatusOK, devices)
}

// GET /api/devices/{mac}
func (g *goWeb) getDevice(w http.ResponseWriter, r *http.Request) {
	if dev := g.findDevice(r); dev != nil {
		g.responseJSON(w, http.StatusOK, MakeDeviceJSON(dev))
	} else {
		g.statusJSON(w, http.StatusNotFound, "Unknown device.")
	}
}

// POST /api/devices/{mac}/upgrade
func (g *goWeb) upgradeDevice(w http.ResponseWriter, r *http.Request) {
	dev := g.findDevice(r)
	if dev == nil {
		g.statusJSON(w, http.StatusNotFound, "Unknown device.")
		return
	}
	if err := dev.Upgrade(); err != nil {
		g.statusJSON(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}
	g.statusJSON(w, http.StatusOK, "Firmware upgrade for %s enqueued.", dev.MacAddress)
}

// POST /api/devices/{mac}/provision
func (g *goWeb) provisionDevice(w http.ResponseWriter, r *http.Request) {
	dev := g.findDevice(r)
	if dev == nil {
		g.statusJSON(w, http.StatusNotFound, "Unknown device.")
		return
	}
	if err := dev.Provision(); err != nil {
		g.statusJSON(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}
	g.statusJSON(w, http.StatusOK, "Provisioning for %s enqueued.", dev.MacAddress)
}

// POST /api/devices/{mac}/reboot
func (g *goWeb) rebootDevice(w http.ResponseWriter, r *http.Request) {
	dev := g.findDevice(r)
	if dev == nil {
		g.statusJSON(w, http.StatusNotFound, "Unknown device.")
		return
	}
	if err := dev.Reboot(); err != nil {
		g.statusJSON(w, http.StatusUnprocessableEntity, "%v", err)
		return
	}
	g.statusJSON(w, http.StatusOK, "Rebooting device %s.", dev.MacAddress)
}

func (g *goWeb) findDevice(r *http.Request) *provisioner.Device {
	vars := mux.Vars(r)
	if mac, found := vars["mac"]; found {
		return g.config.FindDevice(mac)
	}
	return nil
}
