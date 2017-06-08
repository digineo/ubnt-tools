package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/digineo/ubnt-go/discovery"
)

var syslog = flag.Bool("syslog", false, "Disable log timestamps and redirect output to stdout")

func main() {
	flag.Parse()

	if *syslog {
		log.SetFlags(0)
		log.SetOutput(os.Stdout)
	}

	interfaces := []string{}
	for _, iface := range flag.Args() {
		interfaces = append(interfaces, iface)
	}

	notifier, err := discovery.AutoDiscover(interfaces...)
	if err != nil {
		log.Fatal(err)
	}
	go logDevice(notifier)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	os.Exit(0)
}

func logDevice(n <-chan *discovery.Device) {
	for device := range n {
		log.Printf("[discovery] found new device:\n%s", device)
	}
}
