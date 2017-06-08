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
	discoveryPort      = uint16(10001)
	discoveryBroadcast = "255.255.255.255"
	discoveryMulticast = "233.89.188.1"
)

var (
	udpConnections []*net.UDPConn
	started        = false
	discovered     = make(map[string]*Device)
	helloPacket    = map[int][]byte{
		1: {0x01, 0x00, 0x00, 0x00},
		2: {0x02, 0x0a, 0x00, 0x00},
	}
	mutex = &sync.RWMutex{}
)
var incoming = make(chan *Packet, 32)

// AutoDiscover starts the UBNT auto discovery mechanism. It returns a
// notifier channel which receives newly discovered devices (i.e. a
// device already seen won't be send again).
func AutoDiscover(interfaceNames ...string) (<-chan *Device, error) {
	if started {
		return nil, fmt.Errorf("already initialized")
	}
	started = true
	var locals []*net.UDPAddr

	for _, interfaceName := range interfaceNames {
		for _, addr := range interfaceAddresses(interfaceName) {
			locals = append(locals, &net.UDPAddr{IP: addr})
		}
	}

	if len(locals) == 0 {
		return nil, fmt.Errorf("no local addresses on interface %v found", interfaceNames)
	}

	errs := listenMulticast(locals)
	if len(errs) > 0 {
		for i, e := range errs {
			log.Printf("Error %d: %s", i, e.Error())
		}
		return nil, fmt.Errorf("Got errors, see log for details")
	}

	notifier := make(chan *Device, 16)
	go pingDevices()
	go handleIncoming(notifier)
	return notifier, nil
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
func pingDevices() {
	var (
		duration = 4 * time.Second
		max      = 15 * time.Second
		backoff  = 1.02
	)

	for {
		pingMulticast(discoveryMulticast, discoveryPort, helloPacket[1])
		pingMulticast(discoveryBroadcast, discoveryPort, helloPacket[1])

		log.Printf("[discovery] sent broadcast, will send again in %v", duration)
		time.Sleep(duration)
		if duration < max {
			duration = time.Duration(backoff * float64(duration))
		}
	}
}

func pingMulticast(addr string, port uint16, msg []byte) {
	udpAddr := &net.UDPAddr{
		IP:   net.ParseIP(addr),
		Port: int(port),
	}

	for _, conn := range udpConnections {
		conn.WriteToUDP(msg, udpAddr)
	}
}

func listenMulticast(addrs []*net.UDPAddr) (errs []error) {
	for _, addr := range addrs {
		if conn, err := net.ListenUDP("udp", addr); err != nil {
			errs = append(errs, err)
		} else {
			udpConnections = append(udpConnections, conn)
		}
	}
	if len(errs) == 0 {
		for _, conn := range udpConnections {
			log.Printf("[discovery] listen on %s", conn.LocalAddr())
			go packetHandler(conn)
		}
	} else {
		for _, conn := range udpConnections {
			conn.Close()
		}
	}
	return errs
}

func packetHandler(conn *net.UDPConn) {
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
			incoming <- packet
		} else {
			log.Printf("Could not parse packet: %s\n%s\n", err.Error(), hex.Dump(cpy[:n]))
		}
	}
}

func handleIncoming(notifier chan<- *Device) {
	for packet := range incoming {
		dev := packet.Device()
		mutex.RLock()
		old, seen := discovered[dev.MacAddress]
		mutex.RUnlock()

		if !seen || !old.RecentlySeen(1*time.Minute) {
			notifier <- dev
		}

		mutex.Lock()
		if seen {
			old.Merge(dev)
		} else {
			discovered[dev.MacAddress] = dev
		}
		mutex.Unlock()
	}
}

// List all discovered devices so far. Will create duplicates of the
// actual device list, so that it'll be safe to work with the result.
func List() (list []*Device) {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, dev := range discovered {
		list = append(list, dev.Clone())
	}
	return
}

// Find will search the list of discovered devices for an entry with
// matching MAC address, and return a duplicate.
func Find(mac string) *Device {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, dev := range discovered {
		if mac == dev.MacAddress {
			return dev.Clone()
		}
	}
	return nil
}
