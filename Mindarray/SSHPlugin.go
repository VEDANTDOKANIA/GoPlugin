package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
	"time"
)

func sshData(JsonData map[string]interface{}, c chan string) {
	sshHost := (JsonData["IP address"]).(string)
	sshPort := int(JsonData["Port"].(float64))
	sshUser := JsonData["User-Name"].(string)
	sshPassword := JsonData["password"].(string)

	config := &ssh.ClientConfig{
		Timeout:         10 * time.Second,
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Config: ssh.Config{Ciphers: []string{
			"aes128-ctr", "aes192-ctr", "aes256-ctr",
		}},
	}
	config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	sshClient, er := ssh.Dial("tcp", addr, config)

	if er != nil {
		fmt.Println(er)
	}
	structuredmap := make(map[string]interface{})
	var option int

	for option = 1; option < 5; option++ {
		session, er := sshClient.NewSession()
		if er != nil {
			fmt.Println(er)
		}
		defer sshClient.Close()

		//fmt.Scan(&option)
		if option == 1 {

			combo, er := session.CombinedOutput("free -m ")

			output := string(combo)

			res1 := strings.Split(output, "\n")
			res2 := strings.Split(res1[1], "          ")
			res3 := strings.Split(res2[1], "        ")

			var totalMemory, _ = strconv.ParseInt(res3[0], 10, 16)
			var usedMemory, _ = strconv.ParseInt(res3[1], 10, 16)
			var percentageused = (usedMemory * 100) / totalMemory
			var percentageleft = 100 - percentageused
			structuredmap["Memory_Used_Percentage"] = percentageused
			structuredmap["Memory_Available_Percentage"] = percentageleft

			//fmt.Println(structuredmap)
			if er != nil {
			}

		}
		if option == 2 {
			combo, er := session.CombinedOutput("df -h | awk '{if ($1 != \"Filesystem\") print $1 \" \" $5}'")
			output := string(combo)
			res1 := strings.Split(output, "\n")
			var maplist []map[string]interface{}
			for i := 0; i < len(res1)-1; i++ {
				m := make(map[string]interface{})
				res2 := strings.Split(res1[i], " ")
				res3 := strings.Split(res2[1], "%")
				value, _ := strconv.ParseInt(res3[0], 10, 16)
				m["Disk.id"] = res2[0]
				m["Disk.UsedPercent"] = value
				m["Disk.AvailablePercent"] = 100 - value
				maplist = append(maplist, m)
			}
			structuredmap["Disk"] = maplist
			//fmt.Println(structuredmap)

			if er != nil {
			}

		}
		if option == 3 {
			combo, er := session.CombinedOutput("ps aux | awk '{if($3 !=\"%CPU\")print $3 \" \" $11\" \"$2}'")
			output := string(combo)
			//fmt.Println(output)
			res1 := strings.Split(output, "\n")
			//fmt.Println(res1)
			var maplist []map[string]interface{}
			for i := 0; i < len(res1)-1; i++ {
				m := make(map[string]interface{})
				res2 := strings.Split(res1[i], " ")
				value, _ := strconv.ParseInt(res2[0], 10, 16) //PercentProcesstime
				pid, _ := strconv.ParseInt(res2[2], 10, 16)
				m["Process.ID"] = pid
				m["Process.Name"] = res2[1]
				m["Process.%ProcessTime"] = value
				maplist = append(maplist, m)
			}
			structuredmap["Process"] = maplist
			//fmt.Println(structuredmap)
			if er != nil {
			}

		}
		if option == 4 {

			combo, er := session.CombinedOutput("mpstat -P ALL\n")
			output := string(combo)
			res1 := strings.Split(output, "\n")
			var maplist []map[string]interface{}
			for i := 4; i < len(res1)-1; i++ {
				m := make(map[string]interface{})
				res2 := strings.Split(res1[i], "   ")
				//user_percentage, _ := strconv.ParseFloat(strings.Trim(res2[2], " "), 10) //user percentage
				//system_percentage, _ := strconv.ParseFloat(strings.Trim(res2[3], " "), 10)
				ideal_percentage, _ := strconv.ParseFloat(res2[11], 10)
				var cpu_utilization = 100.00 - ideal_percentage

				m["CPU.Core"] = res2[1]
				m["CPU.LoadPercent"] = cpu_utilization
				maplist = append(maplist, m)
			}
			structuredmap["CPU"] = maplist
			//fmt.Println(structuredmap)
			if er != nil {
			}

		}

	}
	structuredmap["IP-Address"] = sshHost
	result, _ := json.Marshal(structuredmap)

	c <- string(result)

}
