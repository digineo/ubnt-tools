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

	discover, err := discovery.AutoDiscover(logDevice, interfaces...)
	if err != nil {
		log.Fatal(err)
	}
	defer discover.Close()

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}

func logDevice(device *discovery.Device) {
	log.Printf("[discovery] found new device:\n%s", device)
}
