package snmp

import (
	"net"
	"strconv"
)

type MIB struct {
	UniFi ubntUniFi // .1.3.6.1.4.1.41112.1.6
}

func (obj *MIB) Assign(pdu SnmpPDU) {
	switch pdu.Name {
	case "6":
		obj.UniFi.Assign(pdu)
	}
}

type UniFi struct {
	Wireless unifiApWireless // .1.3.6.1.4.1.41112.1.6.1
	If       unifiApIf       // .1.3.6.1.4.1.41112.1.6.2
	System   unifiApSystem   // .1.3.6.1.4.1.41112.1.6.3
}

func (obj *UniFi) Assign(pdu SnmpPDU) {
	switch pdu.Name {
	case "1":
		obj.Wireless.Assign(pdu)
	case "2":
		obj.If.Assign(pdu)
	case "3":
		obj.System.Assign(pdu)
	}
}

type ApIf struct {
	IfTable []IfEntry // .1.3.6.1.4.1.41112.1.6.2.1
}

func (obj *ApIf) Assign(pdu SnmpPDU) {
	switch pdu.Name {
	case "1":
		//[]IfEntry
		i := LastIndexByte(pdu.Name, ".")
		index, _ := strconv.Atoi(pdu.Name[i+1:])
		if len(obj.IfTable) == index {
			obj.IfTable = append(obj.IfTable, IfEntry{})
		}
		if len(obj.IfTable) > index {
			pdu.Name = pdu.Name[:i]
			obj.IfTable[index].Assign(pdu)
		}
	}
}

type ApSystem struct {
	IP       net.IP // .1.3.6.1.4.1.41112.1.6.3.1
	Isolated int    // .1.3.6.1.4.1.41112.1.6.3.2
	Model    string // .1.3.6.1.4.1.41112.1.6.3.3
	Uplink   string // .1.3.6.1.4.1.41112.1.6.3.4
	Uptime   uint32 // .1.3.6.1.4.1.41112.1.6.3.5
	Version  string // .1.3.6.1.4.1.41112.1.6.3.6
}

func (obj *ApSystem) Assign(pdu SnmpPDU) {
	switch pdu.Name {
	case "1":
		obj.IP = pduIP(pdu)
	case "2":
		obj.Isolated = pduInt(pdu)
	case "3":
		obj.Model = pduString(pdu)
	case "4":
		obj.Uplink = pduString(pdu)
	case "5":
		obj.Uptime = pduUint32(pdu)
	case "6":
		obj.Version = pduString(pdu)
	}
}

type ApWireless struct {
	RadioTable []RadioEntry // .1.3.6.1.4.1.41112.1.6.1.1
	VapTable   []VapEntry   // .1.3.6.1.4.1.41112.1.6.1.2
}

func (obj *ApWireless) Assign(pdu SnmpPDU) {
	switch pdu.Name {
	case "1":
		//[]RadioEntry
		i := LastIndexByte(pdu.Name, ".")
		index, _ := strconv.Atoi(pdu.Name[i+1:])
		if len(obj.RadioTable) == index {
			obj.RadioTable = append(obj.RadioTable, RadioEntry{})
		}
		if len(obj.RadioTable) > index {
			pdu.Name = pdu.Name[:i]
			obj.RadioTable[index].Assign(pdu)
		}
	case "2":
		//[]VapEntry
		i := LastIndexByte(pdu.Name, ".")
		index, _ := strconv.Atoi(pdu.Name[i+1:])
		if len(obj.VapTable) == index {
			obj.VapTable = append(obj.VapTable, VapEntry{})
		}
		if len(obj.VapTable) > index {
			pdu.Name = pdu.Name[:i]
			obj.VapTable[index].Assign(pdu)
		}
	}
}

type IfEntry struct {
	FullDuplex  int    // .1.3.6.1.4.1.41112.1.6.2.1.1.2
	IP          net.IP // .1.3.6.1.4.1.41112.1.6.2.1.1.3
	Mac         []byte // .1.3.6.1.4.1.41112.1.6.2.1.1.4
	Name        string // .1.3.6.1.4.1.41112.1.6.2.1.1.5
	RxBytes     uint32 // .1.3.6.1.4.1.41112.1.6.2.1.1.6
	RxDropped   uint32 // .1.3.6.1.4.1.41112.1.6.2.1.1.7
	RxError     uint32 // .1.3.6.1.4.1.41112.1.6.2.1.1.8
	RxMulticast uint32 // .1.3.6.1.4.1.41112.1.6.2.1.1.9
	RxPackets   uint32 // .1.3.6.1.4.1.41112.1.6.2.1.1.10
	Speed       int32  // .1.3.6.1.4.1.41112.1.6.2.1.1.11
	TxBytes     uint32 // .1.3.6.1.4.1.41112.1.6.2.1.1.12
	TxDropped   uint32 // .1.3.6.1.4.1.41112.1.6.2.1.1.13
	TxError     uint32 // .1.3.6.1.4.1.41112.1.6.2.1.1.14
	TxPackets   uint32 // .1.3.6.1.4.1.41112.1.6.2.1.1.15
	Up          int    // .1.3.6.1.4.1.41112.1.6.2.1.1.16
}

func (obj *IfEntry) Assign(pdu SnmpPDU) {
	switch pdu.Name {
	case "2":
		obj.FullDuplex = pduInt(pdu)
	case "3":
		obj.IP = pduIP(pdu)
	case "4":
		obj.Mac = pduBytes(pdu)
	case "5":
		obj.Name = pduString(pdu)
	case "6":
		obj.RxBytes = pduUint32(pdu)
	case "7":
		obj.RxDropped = pduUint32(pdu)
	case "8":
		obj.RxError = pduUint32(pdu)
	case "9":
		obj.RxMulticast = pduUint32(pdu)
	case "10":
		obj.RxPackets = pduUint32(pdu)
	case "11":
		obj.Speed = pduInt32(pdu)
	case "12":
		obj.TxBytes = pduUint32(pdu)
	case "13":
		obj.TxDropped = pduUint32(pdu)
	case "14":
		obj.TxError = pduUint32(pdu)
	case "15":
		obj.TxPackets = pduUint32(pdu)
	case "16":
		obj.Up = pduInt(pdu)
	}
}

type RadioEntry struct {
	Name      string // .1.3.6.1.4.1.41112.1.6.1.1.1.2
	Radio     string // .1.3.6.1.4.1.41112.1.6.1.1.1.3
	RxPackets uint32 // .1.3.6.1.4.1.41112.1.6.1.1.1.4
	TxPackets uint32 // .1.3.6.1.4.1.41112.1.6.1.1.1.5
	CuTotal   int32  // .1.3.6.1.4.1.41112.1.6.1.1.1.6
	CuSelfRx  int32  // .1.3.6.1.4.1.41112.1.6.1.1.1.7
	CuSelfTx  int32  // .1.3.6.1.4.1.41112.1.6.1.1.1.8
	OtherBss  int32  // .1.3.6.1.4.1.41112.1.6.1.1.1.9
}

func (obj *RadioEntry) Assign(pdu SnmpPDU) {
	switch pdu.Name {
	case "2":
		obj.Name = pduString(pdu)
	case "3":
		obj.Radio = pduString(pdu)
	case "4":
		obj.RxPackets = pduUint32(pdu)
	case "5":
		obj.TxPackets = pduUint32(pdu)
	case "6":
		obj.CuTotal = pduInt32(pdu)
	case "7":
		obj.CuSelfRx = pduInt32(pdu)
	case "8":
		obj.CuSelfTx = pduInt32(pdu)
	case "9":
		obj.OtherBss = pduInt32(pdu)
	}
}

type VapEntry struct {
	BssID       []byte // .1.3.6.1.4.1.41112.1.6.1.2.1.2
	Ccq         int32  // .1.3.6.1.4.1.41112.1.6.1.2.1.3
	Channel     int32  // .1.3.6.1.4.1.41112.1.6.1.2.1.4
	ExtChannel  int32  // .1.3.6.1.4.1.41112.1.6.1.2.1.5
	EssID       string // .1.3.6.1.4.1.41112.1.6.1.2.1.6
	Name        string // .1.3.6.1.4.1.41112.1.6.1.2.1.7
	NumStations int32  // .1.3.6.1.4.1.41112.1.6.1.2.1.8
	Radio       string // .1.3.6.1.4.1.41112.1.6.1.2.1.9
	RxBytes     uint32 // .1.3.6.1.4.1.41112.1.6.1.2.1.10
	RxCrypts    uint32 // .1.3.6.1.4.1.41112.1.6.1.2.1.11
	RxDropped   uint32 // .1.3.6.1.4.1.41112.1.6.1.2.1.12
	RxErrors    uint32 // .1.3.6.1.4.1.41112.1.6.1.2.1.13
	RxFrags     uint32 // .1.3.6.1.4.1.41112.1.6.1.2.1.14
	RxPackets   uint32 // .1.3.6.1.4.1.41112.1.6.1.2.1.15
	TxBytes     uint32 // .1.3.6.1.4.1.41112.1.6.1.2.1.16
	TxDropped   uint32 // .1.3.6.1.4.1.41112.1.6.1.2.1.17
	TxErrors    uint32 // .1.3.6.1.4.1.41112.1.6.1.2.1.18
	TxPackets   uint32 // .1.3.6.1.4.1.41112.1.6.1.2.1.19
	TxRetries   uint32 // .1.3.6.1.4.1.41112.1.6.1.2.1.20
	TxPower     int32  // .1.3.6.1.4.1.41112.1.6.1.2.1.21
	Up          int    // .1.3.6.1.4.1.41112.1.6.1.2.1.22
	Usage       string // .1.3.6.1.4.1.41112.1.6.1.2.1.23
}

func (obj *VapEntry) Assign(pdu SnmpPDU) {
	switch pdu.Name {
	case "2":
		obj.BssID = pduBytes(pdu)
	case "3":
		obj.Ccq = pduInt32(pdu)
	case "4":
		obj.Channel = pduInt32(pdu)
	case "5":
		obj.ExtChannel = pduInt32(pdu)
	case "6":
		obj.EssID = pduString(pdu)
	case "7":
		obj.Name = pduString(pdu)
	case "8":
		obj.NumStations = pduInt32(pdu)
	case "9":
		obj.Radio = pduString(pdu)
	case "10":
		obj.RxBytes = pduUint32(pdu)
	case "11":
		obj.RxCrypts = pduUint32(pdu)
	case "12":
		obj.RxDropped = pduUint32(pdu)
	case "13":
		obj.RxErrors = pduUint32(pdu)
	case "14":
		obj.RxFrags = pduUint32(pdu)
	case "15":
		obj.RxPackets = pduUint32(pdu)
	case "16":
		obj.TxBytes = pduUint32(pdu)
	case "17":
		obj.TxDropped = pduUint32(pdu)
	case "18":
		obj.TxErrors = pduUint32(pdu)
	case "19":
		obj.TxPackets = pduUint32(pdu)
	case "20":
		obj.TxRetries = pduUint32(pdu)
	case "21":
		obj.TxPower = pduInt32(pdu)
	case "22":
		obj.Up = pduInt(pdu)
	case "23":
		obj.Usage = pduString(pdu)
	}
}
