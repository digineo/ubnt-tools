package provisioner

import (
	"bufio"
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/digineo/ubnt-tools/discovery"
	pssh "github.com/digineo/ubnt-tools/provisioner/ssh"
	"golang.org/x/crypto/ssh"
)

// Device is a wrapper around discovery.Device, and annotates a primary
// IP address, the latest available Firmware and/or system config.
//
// It also knows how to communicate with the device.
type Device struct {
	*discovery.Device
	IPAddress        string
	firmwarePath     string
	systemConfigPath string
	RebootedAt       time.Time

	authMethods []ssh.AuthMethod

	busy    bool
	busyMsg string
	busyMtx sync.RWMutex
}

// CanUpgrade indicates, whether new firmware image is available
func (d *Device) CanUpgrade() bool {
	return d.firmwarePath != ""
}

// HasConfig indicates, whether system config is available
func (d *Device) HasConfig() bool {
	return d.systemConfigPath != ""
}

// IsBusy states whether or not this Device is ready to receive commands.
func (d *Device) IsBusy() bool {
	d.busyMtx.RLock()
	defer d.busyMtx.RUnlock()
	return d.busy
}

func (d *Device) setBusy(msg string) error {
	d.busyMtx.Lock()
	defer d.busyMtx.Unlock()
	if d.busy {
		return fmt.Errorf("Device is busy (%s)", d.busyMsg)
	}
	d.busy = true
	d.busyMsg = msg
	return nil
}

// Status gives a human-readable status information about this device. The
// status may be "idle", "upgrading", "provisioning", or "rebooting". Note
// that this status text only indicates a current event, when this device
// is actually marked busy. Otherwise, you'll get the _last_ state.
func (d *Device) Status() string {
	if d.IsBusy() {
		return d.busyMsg
	}
	if d.RebootedAt.After(d.LastSeenAt) {
		return "rebooting"
	}
	return "idle"
}

// Provision updates the system config on the remote device.
func (d *Device) Provision() error {
	if !d.HasConfig() {
		return fmt.Errorf("No device configuration found for %s", d.MacAddress)
	}

	return d.withSSHClient("provisioning", d.doProvision)
}

// runs in background-goroutine
func (d *Device) doProvision(c *ssh.Client) {
	d.log("Start provisioning...")
	var sessionError error

	remotePath := "/tmp/system.cfg"
	if sessionError = pssh.UploadFile(c, d.systemConfigPath, remotePath); sessionError != nil {
		d.log("Upload failed: %v", sessionError)
		return
	}
	d.log("local(%s) -> remote(%s) 100%%", d.systemConfigPath, remotePath)

	if _, sessionError = pssh.ExecuteCommand(c, "/usr/bin/cfgmtd -w -p /etc/"); sessionError != nil {
		d.log("Could not save configuration: %v", sessionError)
		return
	}
	d.log("Configuration saved")

	if _, sessionError = pssh.ExecuteCommand(c, "/usr/bin/reboot"); sessionError != nil {
		d.log("Reboot failed: %v", sessionError)
		return
	}
	d.markReboot(5 * time.Second)
	d.log("Reboot succeeded")
}

// Upgrade uploads the firmware image to the remote device and starts
// the upgrade process.
func (d *Device) Upgrade() error {
	if !d.CanUpgrade() {
		return fmt.Errorf("cannot safely upgrade device %s", d.MacAddress)
	}
	return d.withSSHClient("upgrading", d.doUpgrade)
}

func (d *Device) doUpgrade(c *ssh.Client) {
	d.log("Start upgrading...")
	var sessionError error

	remotePath := "/tmp/fwupdate.bin"
	if sessionError = pssh.UploadFile(c, d.firmwarePath, remotePath); sessionError != nil {
		d.log("Upload failed: %v", sessionError)
		return
	}
	d.log("local(%s) -> remote(%s) 100%%", d.firmwarePath, remotePath)

	_, sessionError = pssh.ExecuteCommand(c, "/usr/bin/ubntbox fwupdate.real -c "+remotePath)
	if sessionError != nil {
		d.log("Firmware check failed: %v", sessionError)
		return
	}
	d.log("Firmware check succeeded")

	sessionError = pssh.WithinSession(c, func(s *ssh.Session) error {
		reader, err := s.StderrPipe()
		if err != nil {
			return err
		}

		if err := s.Start("/usr/bin/ubntbox fwupdate.real -m " + remotePath); err != nil {
			return err
		}

		bufreader := bufio.NewReader(reader)
		for {
			line, err := bufreader.ReadString('\n')
			d.log("remote: |%s|", strings.TrimRight(line, "\n"))
			if err != nil {
				return err
			}

			if strings.HasPrefix(line, "Done") {
				time.Sleep(15 * time.Second)
				s.Close()
				return nil
			}
		}
	})

	if sessionError != nil {
		d.log("Could not upgrade firmware: %v", sessionError)
		return
	}
	d.markReboot(30 * time.Second)
	d.log("Firmware upgrade succeeded")
}

// Reboot issues a reboot on the device.
func (d *Device) Reboot() error {
	return d.withSSHClient("rebooting", func(c *ssh.Client) {
		if _, sessionError := pssh.ExecuteCommand(c, "/usr/bin/reboot"); sessionError != nil {
			d.log("Reboot failed: %v", sessionError)
			return
		}
		d.markReboot(5 * time.Second)
		d.log("Reboot succeeded")
	})
}

func (d *Device) withSSHClient(msg string, callback func(*ssh.Client)) error {
	if err := d.setBusy(msg); err != nil {
		return err
	}

	client := d.getSSHClient()
	if client == nil {
		d.busy = false
		return fmt.Errorf("Could not obtain SSH client")
	}
	d.log("Got a client")

	go func() {
		callback(client)
		d.log("Callback succeeded")
		d.busy = false
		client.Close()
	}()

	return nil
}

func (d *Device) getSSHClient() *ssh.Client {
	clientConfig := &ssh.ClientConfig{
		Timeout:         2 * time.Second,
		User:            "ubnt",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	for i, m := range d.authMethods {
		clientConfig.Auth = []ssh.AuthMethod{m}
		authType := reflect.TypeOf(m).String()

		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", d.IPAddress), clientConfig)
		if err != nil {
			d.log("(try %d) %s authentication failed with %v", i+1, authType, err)
			continue
		}

		d.log("(try %d) %s authentication succeeded", i+1, authType)
		return client
	}

	return nil
}

// Prints a log message prefixed with "[Device aa:bb:cc:dd:ee:ff]".
func (d *Device) log(message string, v ...interface{}) {
	log.Printf(fmt.Sprintf("[Device %s] %s", d.MacAddress, message), v...)
}

// markReboot sets the RebootedAt flat to a time in the future. This is
// used to detect reboot cycles, which may not be effective immediately,
// and hence makes the device misleadingly available/idle in the UI.
func (d *Device) markReboot(inFuture time.Duration) {
	d.RebootedAt = time.Now().Add(inFuture)
}
