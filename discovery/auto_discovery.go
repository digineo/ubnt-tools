package discovery

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	discoveryPort      = 10001
	discoveryBroadcast = "255.255.255.255"
	discoveryMulticast = "233.89.188.1"

	backoff     = 1.02
	minDuration = 4 * time.Second
	maxDuration = 15 * time.Second
)

var (
	helloPacket = map[int][]byte{
		1: {0x01, 0x00, 0x00, 0x00},
		2: {0x02, 0x0a, 0x00, 0x00},
	}
)

type NotifyHandler func(*Device)

type Discover struct {
	NotifyHandler NotifyHandler
	connections   []*net.UDPConn
	incoming      chan *Packet
	stop          chan interface{}
	devices       map[string]*Device // discovered devices
	mutex         sync.RWMutex
	wg            sync.WaitGroup
}

// AutoDiscover starts the UBNT auto discovery mechanism. It returns a
// notifier channel which receives newly discovered devices (i.e. a
// device already seen won't be send again). You can stop the discovery
// by closing the quit channel.
func AutoDiscover(notify NotifyHandler, interfaceNames ...string) (d *Discover, err error) {
	var locals []*net.UDPAddr

	for _, interfaceName := range interfaceNames {
		for _, addr := range interfaceAddresses(interfaceName) {
			locals = append(locals, &net.UDPAddr{IP: addr})
		}
	}

	if len(locals) == 0 {
		err = fmt.Errorf("no local addresses on interface %v found", interfaceNames)
		return
	}

	d = &Discover{
		devices:       make(map[string]*Device),
		stop:          make(chan interface{}),
		incoming:      make(chan *Packet, 32),
		NotifyHandler: notify,
	}

	errs := d.listenMulticast(locals)
	if len(errs) > 0 {
		for i, e := range errs {
			log.Printf("Error %d: %s", i, e.Error())
		}
		err = fmt.Errorf("Got errors, see log for details")
		return
	}

	go d.pingDevices()
	go d.handleIncoming()

	return d, nil
}

// Close closes ...
func (d *Discover) Close() {
	close(d.stop)
	for _, conn := range d.connections {
		conn.Close()
	}
	d.wg.Wait()
	close(d.incoming)
}

func interfaceAddresses(ifaceName string) (result []net.IP) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		log.Fatal(err)
	}

	if iface.Flags&(net.FlagMulticast|net.FlagBroadcast) == 0 {
		// skip non-multicast/non-broadcast interfaces
		log.Fatalf("interface %s has no multicast/broadcast flags", ifaceName)
	}

	addresses, err := iface.Addrs()
	if err != nil {
		log.Fatal(err)
	}

	for _, addr := range addresses {
		ipnet, ok := addr.(*net.IPNet)
		if ok && ipnet.IP.To4() != nil {
			result = append(result, ipnet.IP)
		}
	}
	return
}

// pings devices and sleeps with exponential back-off (up to a maximum)
func (d *Discover) pingDevices() {
	var duration time.Duration

	for {
		select {
		case <-d.stop:
			return
		case <-time.After(duration):
			d.pingMulticast(discoveryMulticast, helloPacket[1])
			d.pingMulticast(discoveryBroadcast, helloPacket[1])

			if duration.Nanoseconds() == 0 {
				duration = minDuration
			} else if duration < maxDuration {
				duration = time.Duration(backoff * float64(duration))
			}
			log.Printf("[discovery] sent broadcast, will send again in %v", duration)
		}
	}
}

func (d *Discover) pingMulticast(addr string, msg []byte) {
	udpAddr := &net.UDPAddr{
		IP:   net.ParseIP(addr),
		Port: discoveryPort,
	}

	for _, conn := range d.connections {
		conn.WriteToUDP(msg, udpAddr)
	}
}

func (d *Discover) listenMulticast(addrs []*net.UDPAddr) (errs []error) {
	for _, addr := range addrs {
		if conn, err := net.ListenUDP("udp", addr); err != nil {
			errs = append(errs, err)
		} else {
			d.connections = append(d.connections, conn)
		}
	}
	if len(errs) == 0 {
		for _, conn := range d.connections {
			log.Printf("[discovery] listen on %s", conn.LocalAddr())
			go d.packetHandler(conn)
		}
	} else {
		for _, conn := range d.connections {
			conn.Close()
		}
	}
	return errs
}

func (d *Discover) packetHandler(conn *net.UDPConn) {
	d.wg.Add(1)
	defer d.wg.Done()

	buf := make([]byte, 1500)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Print(err)
			return
		}

		if n <= 4 {
			continue // cannot possibly be a discovery response
		}

		cpy := make([]byte, 1500)
		copy(cpy, buf[:n])

		if packet, err := ParsePacket(cpy[:n]); err == nil {
			d.incoming <- packet
		} else {
			log.Printf("Could not parse packet: %s\n%s\n", err.Error(), hex.Dump(cpy[:n]))
		}
	}
}

func (d *Discover) handleIncoming() {
	for packet := range d.incoming {
		dev := packet.Device()
		d.mutex.RLock()
		old, seen := d.devices[dev.MacAddress]
		d.mutex.RUnlock()

		if !seen || !old.RecentlySeen(1*time.Minute) {
			if handler := d.NotifyHandler; d != nil {
				handler(dev)
			}
		}

		d.mutex.Lock()
		if seen {
			old.Merge(dev)
		} else {
			d.devices[dev.MacAddress] = dev
		}
		d.mutex.Unlock()
	}
}

// List all discovered devices so far. Will create duplicates of the
// actual device list, so that it'll be safe to work with the result.
func (d *Discover) List() (list []*Device) {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	for _, dev := range d.devices {
		list = append(list, dev.Clone())
	}
	return
}

// Find will search the list of discovered devices for an entry with
// matching MAC address, and return a duplicate.
func (d *Discover) Find(mac string) *Device {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	for _, dev := range d.devices {
		if mac == dev.MacAddress {
			return dev.Clone()
		}
	}
	return nil
}
