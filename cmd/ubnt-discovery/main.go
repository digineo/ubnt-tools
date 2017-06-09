package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/digineo/goldflags"
	"github.com/digineo/ubnt-tools/discovery"
)

const appName = "ubnt-discovery"

var syslog = flag.Bool("syslog", false, "Disable log timestamps and redirect output to stdout")

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, goldflags.Banner(appName))
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *syslog {
		log.SetFlags(0)
		log.SetOutput(os.Stdout)
	}

	log.Println(goldflags.Banner("ubnt-discovery"))

	interfaces := []string{}
	for _, iface := range flag.Args() {
		interfaces = append(interfaces, iface)
	}

	notifier, quitter, err := discovery.AutoDiscover(interfaces...)
	if err != nil {
		log.Fatal(err)
	}
	go logDevice(notifier)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	close(quitter)
	os.Exit(0)
}

func logDevice(n <-chan *discovery.Device) {
	for device := range n {
		log.Printf("[discovery] found new device:\n%s", device)
	}
}
