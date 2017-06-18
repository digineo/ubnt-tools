package snmp

import "net"

type MIB struct {
	UniFi *UniFi `oid:".1.3.6.1.4.1.41112.1.6"`
}
type UniFi struct {
	Wireless *ApWireless `oid:".1.3.6.1.4.1.41112.1.6.1"`
	If       *ApIf       `oid:".1.3.6.1.4.1.41112.1.6.2"`
	System   *ApSystem   `oid:".1.3.6.1.4.1.41112.1.6.3"`
}
type ApWireless struct {
	RadioTable []RadioEntry `oid:".1.3.6.1.4.1.41112.1.6.1.1"`
	VapTable   []VapEntry   `oid:".1.3.6.1.4.1.41112.1.6.1.2"`
}
type ApIf struct {
	IfTable []IfEntry `oid:".1.3.6.1.4.1.41112.1.6.2.1"`
}
type ApSystem struct {
	IP       net.IP `oid:".1.3.6.1.4.1.41112.1.6.3.1"`
	Isolated int    `oid:".1.3.6.1.4.1.41112.1.6.3.2"`
	Model    string `oid:".1.3.6.1.4.1.41112.1.6.3.3"`
	Uplink   string `oid:".1.3.6.1.4.1.41112.1.6.3.4"`
	Uptime   uint32 `oid:".1.3.6.1.4.1.41112.1.6.3.5"`
	Version  string `oid:".1.3.6.1.4.1.41112.1.6.3.6"`
}
type RadioEntry struct {
	Name      string `oid:".1.3.6.1.4.1.41112.1.6.1.1.1.2"`
	Radio     string `oid:".1.3.6.1.4.1.41112.1.6.1.1.1.3"`
	RxPackets uint32 `oid:".1.3.6.1.4.1.41112.1.6.1.1.1.4"`
	TxPackets uint32 `oid:".1.3.6.1.4.1.41112.1.6.1.1.1.5"`
	CuTotal   int32  `oid:".1.3.6.1.4.1.41112.1.6.1.1.1.6"`
	CuSelfRx  int32  `oid:".1.3.6.1.4.1.41112.1.6.1.1.1.7"`
	CuSelfTx  int32  `oid:".1.3.6.1.4.1.41112.1.6.1.1.1.8"`
	OtherBss  int32  `oid:".1.3.6.1.4.1.41112.1.6.1.1.1.9"`
}
type IfEntry struct {
	FullDuplex  int    `oid:".1.3.6.1.4.1.41112.1.6.2.1.1.2"`
	IP          net.IP `oid:".1.3.6.1.4.1.41112.1.6.2.1.1.3"`
	Mac         []byte `oid:".1.3.6.1.4.1.41112.1.6.2.1.1.4"`
	Name        string `oid:".1.3.6.1.4.1.41112.1.6.2.1.1.5"`
	RxBytes     uint32 `oid:".1.3.6.1.4.1.41112.1.6.2.1.1.6"`
	RxDropped   uint32 `oid:".1.3.6.1.4.1.41112.1.6.2.1.1.7"`
	RxError     uint32 `oid:".1.3.6.1.4.1.41112.1.6.2.1.1.8"`
	RxMulticast uint32 `oid:".1.3.6.1.4.1.41112.1.6.2.1.1.9"`
	RxPackets   uint32 `oid:".1.3.6.1.4.1.41112.1.6.2.1.1.10"`
	Speed       int32  `oid:".1.3.6.1.4.1.41112.1.6.2.1.1.11"`
	TxBytes     uint32 `oid:".1.3.6.1.4.1.41112.1.6.2.1.1.12"`
	TxDropped   uint32 `oid:".1.3.6.1.4.1.41112.1.6.2.1.1.13"`
	TxError     uint32 `oid:".1.3.6.1.4.1.41112.1.6.2.1.1.14"`
	TxPackets   uint32 `oid:".1.3.6.1.4.1.41112.1.6.2.1.1.15"`
	Up          int    `oid:".1.3.6.1.4.1.41112.1.6.2.1.1.16"`
}
type VapEntry struct {
	BssID       []byte `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.2"`
	Ccq         int32  `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.3"`
	Channel     int32  `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.4"`
	ExtChannel  int32  `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.5"`
	EssID       string `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.6"`
	Name        string `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.7"`
	NumStations int32  `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.8"`
	Radio       string `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.9"`
	RxBytes     uint32 `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.10"`
	RxCrypts    uint32 `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.11"`
	RxDropped   uint32 `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.12"`
	RxErrors    uint32 `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.13"`
	RxFrags     uint32 `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.14"`
	RxPackets   uint32 `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.15"`
	TxBytes     uint32 `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.16"`
	TxDropped   uint32 `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.17"`
	TxErrors    uint32 `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.18"`
	TxPackets   uint32 `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.19"`
	TxRetries   uint32 `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.20"`
	TxPower     int32  `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.21"`
	Up          int    `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.22"`
	Usage       string `oid:".1.3.6.1.4.1.41112.1.6.1.2.1.23"`
}
