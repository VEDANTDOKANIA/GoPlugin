package main

import (
	"encoding/json"
	"github.com/masterzen/winrm"
	"math"
	"strconv"
	"strings"
)

func winrmdata(JsonData map[string]interface{}, c chan string) {
	Host := (JsonData["IP address"]).(string)
	Port := int(JsonData["Port"].(float64))
	username := JsonData["User-Name"].(string)
	password := JsonData["password"].(string)

	endpoint := winrm.NewEndpoint(Host, Port, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, username, password)
	if err != nil {
		panic(err)
	}
	structuredmap := make(map[string]interface{})
	var option int
	for option = 1; option < 5; option++ {
		//fmt.Scan(&option)
		if option == 1 {
			a := "aa"
			output := ""
			ac := "Get-WmiObject win32_OperatingSystem |%{\"{0} \n{1} \n{2} \n" +
				"{3}\" -f $_.totalvisiblememorysize, $_.freephysicalmemory, $_.totalvirtualmemorysize, $_.freevirtualmemory}" // Command jo humko run karna hain
			output, _, _, err = client.RunPSWithString(ac, a)
			//fmt.Println(output)
			res1 := strings.Split(output, "\n")
			//fmt.Println(res1[0])
			total_space_memory, _ := strconv.ParseInt(strings.TrimSpace(res1[0]), 10, 64)
			total_space_virtual, _ := strconv.ParseInt(strings.TrimSpace(res1[2]), 10, 64)
			free_space_memory, _ := strconv.ParseInt(strings.TrimSpace(res1[1]), 10, 64)
			free_space_virtual, _ := strconv.ParseInt(strings.TrimSpace(res1[3]), 10, 64)
			total_space := float64(total_space_memory + total_space_virtual)
			free_space := float64(free_space_virtual + free_space_memory)
			percent := float64(free_space/total_space) * 100
			structuredmap["Memory_Used_Percentage"] = percent
			structuredmap["Memory_Available_Percentage"] = 100.0 - percent
		} //Finally Completed with new data structure
		if option == 2 {
			a := "aa"
			output := ""

			ac := "Get-WmiObject win32_logicaldisk |Foreach-Object {$_.DeviceId,$_.Freespace,$_.Size}" // Command jo humko run karna hain
			output, _, _, err = client.RunPSWithString(ac, a)

			res1 := strings.Split(output, "\r\n")

			if len(res1)%3 != 0 {
				res1 = append(res1, string(0), string(0))
			}
			var maplist []map[string]interface{}

			for i := 0; i < len(res1)-1; i++ {
				if i == len(res1) {
					break
				}
				if i%3 == 0 {
					m := make(map[string]interface{})
					drivename := strings.Split(res1[i], ":")
					freespace, _ := strconv.ParseFloat(res1[i+1], 10)
					Size, _ := strconv.ParseFloat(res1[i+2], 10)
					//m[drivename[0]] = make(map[string]float64)
					m["Disk"+".freespace"] = freespace
					m["Disk"+".size"] = Size
					percentage := (freespace / Size) * 100
					if math.IsNaN(percentage) {
						percentage = 0
					}
					m["Disk"+".percentfree"] = percentage
					m["Disk"+".name"] = drivename[0]
					maplist = append(maplist, m)
				}
			}

			structuredmap["Disk"] = maplist
		} //Finally Completed with new data structure
		if option == 3 {
			a := "aa"
			output := ""

			ac := "(Get-Counter -Counter '\\Process(*)\\ID Process','\\Process(*)\\% Processor Time' -ErrorAction SilentlyContinue).counterSamples |Format-List -Property Path,Cookedvalue" // Command jo humko run karna hain
			output, _, _, err = client.RunPSWithString(ac, a)

			//fmt.Println(string(output))

			res1 := strings.Split(output, "\r\n")
			var res2 []string
			var res8 []string
			for i := 0; i < len(res1); i++ {
				res8 = strings.Split(res1[i], ":")
				if len(res8) > 1 {
					res2 = append(res2, res1[i])
				}
			}

			var maplist []map[string]interface{}
			for i := 0; i < len(res2)/2; i = i + 2 {
				datamap := make(map[string]interface{})
				res3 := strings.Split(res2[i], "(")
				res4 := strings.Split(res3[1], ")")
				res5 := strings.Split(res2[i+1], ": ")
				value, _ := strconv.ParseFloat(res5[1], 10)
				datamap["Process.name"] = res4[0]
				datamap["Process.ID"] = value

				j := i + len(res2)/2 + 1
				res6 := strings.Split(res2[j], ": ")
				valueprocess, _ := strconv.ParseFloat(res6[1], 10)
				datamap["Process.%ProcessTime"] = valueprocess
				maplist = append(maplist, datamap)

			}
			structuredmap["Process"] = maplist

		} //FInally Completed with new data structure
		if option == 4 {
			a := "aa"
			output := ""
			ac := "Get-WmiObject win32_processor | select SystemName, LoadPercentage" // Command jo humko run karna hain
			output, _, _, err = client.RunPSWithString(ac, a)
			//fmt.Println(string(output))
			res1 := strings.Split(output, "\r\n")
			var res2 []string
			for i := 0; i < len(res1); i++ {
				if res1[i] != "" {
					res2 = append(res2, res1[i])
				}
			}
			counter := 0
			var maplist []map[string]interface{}
			for i := 2; i < len(res2); i++ {
				datamap := make(map[string]interface{})
				systemName := strings.Split(res2[i], " ")
				//corename := strings.TrimSpace(systemName[0])
				loadvalue, _ := strconv.ParseInt((strings.TrimSpace(systemName[len(systemName)-1])), 10, 64)
				datamap["CPU.core"] = counter
				datamap["CPU.LoadPercent"] = loadvalue
				counter++
				maplist = append(maplist, datamap)
			}
			structuredmap["CPU"] = maplist
		} //Finally Completed with new data structure

	}
	//fmt.Println(structuredmap)
	structuredmap["IP-Address"] = Host

	result, _ := json.Marshal(structuredmap)

	//fmt.Println(string(result))
	c <- string(result)
}
