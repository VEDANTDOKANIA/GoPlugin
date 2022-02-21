package main

import (
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"strconv"
	"time"
)

var interfaceoutput []string

func snmpData(JsonData map[string]interface{}, c chan string) {
	version := (JsonData["Version"]).(string)
	Version1 := g.Version2c
	if version == "v2c" {
		Version1 = g.Version2c
	}
	if version == "v1" {
		Version1 = g.Version1
	}
	if version == "v3" {
		Version1 = g.Version3
	}

	params := &g.GoSNMP{
		Target:    JsonData["IP address"].(string),
		Port:      uint16((JsonData["Port"]).(float64)),
		Community: (JsonData["Community"]).(string),
		Version:   Version1,
		Timeout:   time.Duration(2) * time.Second,
		//Logger:    g.NewLogger(log.New(os.Stdout, "", 0)),
	}
	err1 := params.Connect()
	if err1 != nil {

	}

	defer params.Conn.Close()

	oids := []string{"1.3.6.1.2.1.1.5.0", "1.3.6.1.2.1.1.1.0", "1.3.6.1.2.1.1.6.0", "1.3.6.1.2.1.1.2.0", "1.3.6.1.2.1.1.3.0", "1.3.6.1.2.1.2.1.0"}
	result, _ := params.Get(oids) // Get() accepts up to g.MAX_OIDS

	configmap := make(map[string]interface{})

	for _, variable := range result.Variables {

		switch variable.Name {
		case ".1.3.6.1.2.1.1.5.0":
			configmap["System_Name"] = string(variable.Value.([]byte))
		case ".1.3.6.1.2.1.1.1.0":
			configmap["System_Description"] = string(variable.Value.([]byte))
		case ".1.3.6.1.2.1.1.6.0":
			configmap["System_Loaction"] = string(variable.Value.([]byte))
		case ".1.3.6.1.2.1.1.2.0":
			configmap["System_OID"] = variable.Value
		case ".1.3.6.1.2.1.1.3.0":
			configmap["System_Uptime"] = variable.Value
		case ".1.3.6.1.2.1.2.1.0":
			configmap["Total_Interface"] = variable.Value
		default:
			fmt.Println("Unknown variable Name")
		}

	}

	interfaceCount := (configmap["Total_Interface"]).(int)

	var interfacelist []map[string]interface{}
	interfaceindex := params.BulkWalk("1.3.6.1.2.1.2.2.1.1", WalkFunction)
	if interfaceindex != nil {

	}

	for i := 0; i < interfaceCount; i++ {
		tempmap := make(map[string]interface{})
		tempmap["Interface_Index"], _ = strconv.ParseInt(interfaceoutput[i], 10, 64)
		interfacelist = append(interfacelist, tempmap)
	}
	interfaceoutput = nil

	interfacedescription := params.BulkWalk("1.3.6.1.2.1.2.2.1.2", WalkFunction)
	if interfacedescription != nil {

	}

	for i := 0; i < interfaceCount; i++ {
		tempmap := make(map[string]interface{})
		tempmap = interfacelist[i]
		tempmap["Interface_Description"] = interfaceoutput[i]
		interfacelist[i] = tempmap
	}
	interfaceoutput = nil

	interfaceoutput = nil
	interfacestatus := params.BulkWalk("1.3.6.1.2.1.2.2.1.7", WalkFunction)
	if interfacestatus != nil {

	}

	for i := 0; i < interfaceCount; i++ {
		tempmap := make(map[string]interface{})
		tempmap = interfacelist[i]
		status_value, _ := strconv.ParseInt(interfaceoutput[i], 10, 64)
		if status_value == 1 {
			tempmap["Interface_Status"] = "Up"
			interfacelist[i] = tempmap
		}

		if status_value == 2 {
			tempmap["Interface_Status"] = "Down"
			interfacelist[i] = tempmap
		}
	}
	interfaceoutput = nil
	interfacetype := params.BulkWalk("1.3.6.1.2.1.2.2.1.3", WalkFunction)
	if interfacetype != nil {

	}

	for i := 0; i < interfaceCount; i++ {
		tempmap := make(map[string]interface{})
		tempmap = interfacelist[i]
		type_value, _ := strconv.ParseInt(interfaceoutput[i], 10, 64)
		if type_value == 6 {
			tempmap["Interface_Type"] = "ethernetCsmacd"
			interfacelist[i] = tempmap
		}
		if type_value == 135 {
			tempmap["Interface_Type"] = "l2vlan"
			interfacelist[i] = tempmap
		}
		if type_value == 53 {
			tempmap["Interface_Type"] = "propVirtual"
			interfacelist[i] = tempmap
		}
		if type_value == 24 {
			tempmap["Interface_Type"] = "softwareLoopback"
			interfacelist[i] = tempmap
		}
		if type_value == 131 {
			tempmap["Interface_Type"] = "tunnel "
			interfacelist[i] = tempmap
		}

	}
	interfaceoutput = nil

	configmap["Interface"] = interfacelist

	output, _ := json.Marshal(configmap)

	c <- string(output)

}
func WalkFunction(pdu g.SnmpPDU) error {

	switch pdu.Type {
	case g.OctetString:
		result := pdu.Value.([]byte)
		interfaceoutput = append(interfaceoutput, string(result))
	default:
		result := g.ToBigInt(pdu.Value)
		interfaceoutput = append(interfaceoutput, result.String())
	}

	return nil
}
