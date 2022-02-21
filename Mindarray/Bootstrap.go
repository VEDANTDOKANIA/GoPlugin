package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

func main() {

	decodestring, _ := base64.StdEncoding.DecodeString((os.Args[1]))
	motadatamap := make(map[string][]map[string]interface{})
	var err = json.Unmarshal([]byte(string(decodestring)), &motadatamap)
	if err != nil {
		fmt.Println(err)
	}
	listofmaps := motadatamap["1"]
	c := make(chan string, 1)

	for i := 0; i < len(listofmaps); i++ {
		m2 := listofmaps[i]
		metricType := (m2["Metric_Type"]).(string)
		if metricType == "linux" {
			go sshData(m2, c)
		}
		if metricType == "windows" {
			go winrmdata(m2, c)
		}
		if metricType == "Network_Device" {
			go snmpData(m2, c)
		}
	}
	for i := 0; i < len(listofmaps); i++ {

		fmt.Println(<-c)
		//
		//fmt.Println("New Object")

	}
}
