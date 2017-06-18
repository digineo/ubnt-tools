package snmp

//go:generate ./generate ubntUnifi.csv types.go

import (
	"os"
	"reflect"
	"strings"
	"time"

	"log"

	g "github.com/soniah/gosnmp"
)

func Read(ipAddress string) {

	unifi := &UniFi{}

	snmp := &g.GoSNMP{
		Target:         ipAddress,
		Port:           161,
		Version:        g.Version1,
		Timeout:        time.Duration(3) * time.Second,
		Community:      "public",
		MaxRepetitions: 5,
	}

	err := snmp.Connect()
	if err != nil {
		panic(err)
	}
	defer snmp.Conn.Close()

	cb := func(pdu g.SnmpPDU) error {
		return callback(unifi, pdu)
	}

	val := reflect.ValueOf(unifi).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		oid := typeField.Tag.Get("oid")
		err = snmp.Walk(oid, cb)
		if err != nil {
			log.Printf("Walk Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func callback(object interface{}, pdu g.SnmpPDU) error {
	log.Printf("%s (%d) = %+v", pdu.Name, pdu.Type, pdu.Value)

	val := reflect.ValueOf(object).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		oid := typeField.Tag.Get("oid")
		if pdu.Name == oid || strings.HasPrefix(pdu.Name, oid+".") {
			switch valueField.Kind() {
			case reflect.Ptr:
				log.Println("PTR found", oid, pdu.Name, valueField, typeField.Type, typeField.Name)
				if valueField.IsNil() {
					new := reflect.New(valueField.Type().Elem())
					valueField.Set(new)
				}
				callback(valueField.Interface(), pdu)
			case reflect.Slice:
				log.Println("SLICE", pdu.Name[len(oid):])
			case reflect.String:
				if str, ok := pdu.Value.(string); ok {
					valueField.SetString(str)
				} else {
					log.Println("not string:", pdu)
				}
			default:
				log.Println("DEFAULT", typeField.Type, typeField.Name)
			}
		}
	}
	return nil
}
