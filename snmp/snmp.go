package snmp

//go:generate ./generate ubntUnifi.csv types.go

import (
	g "github.com/soniah/gosnmp"
)

type SnmpPDU g.SnmpPDU
