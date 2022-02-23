package main

import (
	"encoding/json"
	"fmt"
	g "github.com/gosnmp/gosnmp"
	"strconv"
	"strings"
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
	interfaceindex := params.Walk("1.3.6.1.2.1.2.2.1.1", WalkFunction)
	if interfaceindex != nil {

	}
	var interfaceoid []string
	var description string = "1.3.6.1.2.1.2.2.1.2."
	var Interfacetype string = "1.3.6.1.2.1.2.2.1.3."
	var operatingstatus string = "1.3.6.1.2.1.2.2.1.8."
	var adminstatus string = "1.3.6.1.2.1.2.2.1.7."
	var inerrors string = "1.3.6.1.2.1.2.2.1.14."
	var outerrors string = "1.3.6.1.2.1.2.2.1.20."

	for i := 0; i < interfaceCount; i++ {
		Index_value, _ := strconv.ParseInt(interfaceoutput[i], 10, 64)
		descriptionoid := description + strconv.Itoa(int(Index_value))
		interfaceoid = append(interfaceoid, descriptionoid)
		typeoid := Interfacetype + strconv.Itoa(int(Index_value))
		interfaceoid = append(interfaceoid, typeoid)
		operatingstatusoid := operatingstatus + strconv.Itoa(int(Index_value))
		interfaceoid = append(interfaceoid, operatingstatusoid)
		adminstatusoid := adminstatus + strconv.Itoa(int(Index_value))
		interfaceoid = append(interfaceoid, adminstatusoid)
		inerrornoid := inerrors + strconv.Itoa(int(Index_value))
		interfaceoid = append(interfaceoid, inerrornoid)
		outerroroid := outerrors + strconv.Itoa(int(Index_value))
		interfaceoid = append(interfaceoid, outerroroid)
	}
	var interfacelist []map[string]interface{}
	var startIndex = 0
	var endIndex = 60
	var flag = 1
	var inflag = 1
	for flag == 1 {
		if inflag == 0 {
			flag = 0
		}
		InterfaceResult, err := params.Get(interfaceoid[startIndex:endIndex])
		if err != nil {
			fmt.Println(err)
		}
		var count int = 0
		interfaceMap := make(map[string]interface{})
		for _, variable := range InterfaceResult.Variables {

			VariableName := strings.SplitAfter(variable.Name, ".1.3.6.1.2.1.2.2.1.")
			RootOid := VariableName[0] + strings.Split(VariableName[1], ".")[0] + "."
			//VariableName := variable.Name[:len(variable.Name)-1]
			switch RootOid {
			case ".1.3.6.1.2.1.2.2.1.2.":
				descriptionValue := string(variable.Value.([]byte))
				interfaceMap["Interface.Description"] = descriptionValue
				interfaceMap["Interface.Index"], _ = strconv.ParseInt(strings.Split(VariableName[1], ".")[1], 10, 64)
				count++
			case ".1.3.6.1.2.1.2.2.1.3.":
				var typevalue string
				switch (variable.Value).(int) {
				case 6:
					typevalue = "ethernetCsmacd"
				case 1:
					typevalue = "other"
				case 135:
					typevalue = "l2vlan"
				case 53:
					typevalue = "propVirtual"
				case 24:
					typevalue = "softwareLoopback"
				case 131:
					typevalue = "tunnel"
				}
				interfaceMap["Interface.Type"] = typevalue
				count++
			case ".1.3.6.1.2.1.2.2.1.8.":
				var operatingvalue string
				if variable.Value.(int) == 1 {
					operatingvalue = "Up"
				}
				if variable.Value.(int) == 2 {
					operatingvalue = "Down"
				}
				interfaceMap["Interface.Operating_Status"] = operatingvalue
				count++
			case ".1.3.6.1.2.1.2.2.1.7.":
				var Adminvalue string
				if variable.Value.(int) == 1 {
					Adminvalue = "Up"
				}
				if variable.Value.(int) == 2 {
					Adminvalue = "Down"
				}
				interfaceMap["Interface.Admin_Status"] = Adminvalue
				count++
			case ".1.3.6.1.2.1.2.2.1.14.":
				if variable.Value == nil {
					interfaceMap["Interface.InError"] = ""
				} else {
					interfaceMap["Interface.InError"] = variable.Value
				}
				count++
			case ".1.3.6.1.2.1.2.2.1.20.":
				if (variable.Value) == nil {
					interfaceMap["Interface.OutError"] = ""
				} else {
					interfaceMap["Interface.OutError"] = variable.Value
				}

				count++
			}
			if count == 6 {
				interfacelist = append(interfacelist, interfaceMap)
				interfaceMap = make(map[string]interface{})
				count = 0
			}
		}
		//fmt.Println(interfacelist)
		startIndex = endIndex
		endIndex = endIndex + 60
		if endIndex > len(interfaceoid) {
			endIndex = len(interfaceoid)
			inflag = 0
		}
	}
	//fmt.Println(interfacelist)
	configmap["Interface"] = interfacelist
	configmap["IP-Address"] = JsonData["IP address"].(string)

	output, err := json.Marshal(configmap)
	if err != nil {
		fmt.Println(err)
	}
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
