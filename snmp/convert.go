package snmp

import "net"

func pduString(pdu SnmpPDU) string {
	return ""
}
func pduIP(pdu SnmpPDU) net.IP {
	return net.IP{}
}
func pduBytes(pdu SnmpPDU) []byte {
	return []byte{}
}
func pduInt(pdu SnmpPDU) int {
	return -1
}
func pduInt32(pdu SnmpPDU) int32 {
	return -1
}
func pduUint32(pdu SnmpPDU) uint32 {
	return -1
}
