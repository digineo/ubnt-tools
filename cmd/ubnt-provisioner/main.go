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
	"github.com/digineo/ubnt-tools/provisioner"
	"github.com/digineo/ubnt-tools/provisioner/web"
)

const appName = "Digineo Ubnt provisioner"

var (
	configFile = flag.String("c", "./config.yml", "(optional) `path` to config.yml configuration file")
	syslog     = flag.Bool("syslog", false, "Disable log timestamps and redirect output to stdout")

	configuration *provisioner.Configuration
)

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, goldflags.Banner(appName))
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		printExampleConfig()
		flag.PrintDefaults()
	}

	fmt.Println(goldflags.Banner(appName))
	flag.Parse()

	if *syslog {
		log.SetFlags(0)
		log.SetOutput(os.Stdout)
	}

	if configFile == nil || *configFile == "" {
		flag.Usage()
		printExampleConfig()
		os.Exit(1)
	}

	if config, errs := provisioner.LoadConfig(*configFile); len(errs) == 0 {
		configuration = config
	} else {
		log.Printf("Found %d error(s) loading the config file:", len(errs))
		for i, e := range errs {
			log.Printf("Error %d: %s", i, e.Error())
		}
		os.Exit(1)
	}

	discover, _, err := discovery.AutoDiscover(logDevice, configuration.InterfaceNames...)
	if err != nil {
		log.Fatal(err)
	}
	go web.StartWeb(configuration)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	discover.Close()
	os.Exit(0)
}

func logDevice(device *discovery.Device) {
	log.Printf("[main] (re-) discovered %s (%s)", device.Hostname, device.MacAddress)
}

func printExampleConfig() {
	out := os.Stderr
	if *syslog {
		out = os.Stdout
	}
	fmt.Fprintf(out, "Example config file:\n\n%s\n", provisioner.ExampleYAML)
}
