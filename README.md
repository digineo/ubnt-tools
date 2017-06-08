# ubnt-tools

Discovery and provisioning of Ubiquiti AirMax devices.

## `ubnt-discovery`

[![GoDoc](https://godoc.org/github.com/digineo/ubnt-tools/discovery?status.svg)](http://godoc.org/github.com/digineo/ubnt-tools/discovery)

This will install the discovery tool in `$GOPATH/bin/ubnt-discovery`:

    go get github.com/digineo/ubnt-tools/cmd/ubnt-discovery

### Usage

Simply invoke the discovery tool with an interface name:

    ubnt-discovery eth0

This will broadcast the discovery packages (with exponential back-off),
and report back the newly discovered devices:

    2017/06/08 17:39:02 [discovery] listen on 172.16.1.7:49317
    2017/06/08 17:39:02 [discovery] listen on 169.254.0.7:51705
    2017/06/08 17:39:02 [discovery] sent broadcast, will send again in 4s
    2017/06/08 17:39:02 [discovery] found new device:
    Device
      MAC:          80:2a:a8:64:7e:59
      Model:        NanoBeam 5AC 19
      Platform:     N5C
      Firmware:     XC.qca955x.v8.0.2.33352.170327.1907
      Hostname:     NanoBeam 5AC 19
      Booted at:    2017-06-08T16:08:50+02:00
      booted:       1h30m12.00001943s ago
      first seen:   24.031µs ago
      last seen:    24.031µs ago
      IP addresses on interface 80:2a:a8:64:7e:59
        - 192.168.1.20
        - 169.254.126.89
      ESSID:        ubnt
      WMode:        Station
    2017/06/08 17:39:02 [discovery] found new device:
    Device
      MAC:          04:18:d6:83:f8:ec
      Model:
      Platform:     ERLite-3
      Firmware:     EdgeRouter.ER-e100.v1.9.1.4939093.161214.0705
      Hostname:     digineo
      Booted at:    2017-04-10T10:36:11+02:00
      booted:       1423h2m51.000008187s ago
      first seen:   11.759µs ago
      last seen:    11.759µs ago
      IP addresses on interface 04:18:d6:83:f8:ec
        - 172.16.0.1
        - 172.16.2.1
