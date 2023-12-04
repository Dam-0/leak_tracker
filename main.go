package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Leak_nodes struct {
	Leak_nodes []Leak `json:"leaks"`
}

type Leak struct {
	Ip               string    `json:"Ip"`
	Resource_id      string    `json:"resource_id"`
	Leak_count       int       `json:"leak_count"`
	Leak_event_count int       `json:"leak_event_count"`
	Open_ports       []string  `json:"open_ports"`
	Events           []Events  `json:"events"`
	Creation_date    time.Time `json:"creation_date"`
	Update_date      time.Time `json:"update_date"`
}

type Events struct {
	Event_source string    `json:"event_source"`
	Ip           string    `json:"ip"`
	Host         string    `json:"host"`
	Port         string    `json:"port"`
	Protocol     string    `json:"protocol"`
	HTTP         HTTP      `json:"http"`
	Summary      string    `json:"summary"`
	Time         time.Time `json:"time"`
	Network      Network   `json:"network"`
}

type HTTP struct {
	Root string `json:"root"`
	URL  string `json:"url"`
}

type Network struct {
	Organisation   string `json:"organization_name"`
	ASN            int    `json:"asn"`
	Network_subnet string `json:"network"`
}

func main() {
	file, _ := os.ReadFile("report.json")

	data := Leak_nodes{}

	_ = json.Unmarshal([]byte(file), &data)

	for i := 0; i < len(data.Leak_nodes); i++ {
		fmt.Println("IP: ", data.Leak_nodes[i].Ip)
		fmt.Println("Resource ID: ", data.Leak_nodes[i].Resource_id)
		fmt.Println("Leaks Found: ", data.Leak_nodes[i].Leak_count)
		fmt.Println("Total Leaks: ", data.Leak_nodes[i].Leak_event_count)
		fmt.Println("Open Ports: ", data.Leak_nodes[i].Open_ports)

		fmt.Println("\nEVENTS:")

		for x := 0; x < len(data.Leak_nodes[i].Events); x++ {
			fmt.Println("Summary:", data.Leak_nodes[i].Events[x].Summary)
			fmt.Println("Time:", data.Leak_nodes[i].Events[x].Time)
			fmt.Println("Event Source:", data.Leak_nodes[i].Events[x].Event_source)
			fmt.Println("Host IP:", data.Leak_nodes[i].Events[x].Ip)
			fmt.Println("Host Address:", data.Leak_nodes[i].Events[x].Host)
			fmt.Println("Host Port:", data.Leak_nodes[i].Events[x].Port)
			fmt.Println("Host Protocol in Use:", data.Leak_nodes[i].Events[x].Protocol)
			fmt.Println("Root:", data.Leak_nodes[i].Events[x].HTTP.Root)
			fmt.Println("URL:", data.Leak_nodes[i].Events[x].HTTP.URL)
			fmt.Println("Organisation:", data.Leak_nodes[i].Events[x].Network.Organisation)

		}

		fmt.Println("First Seen:", data.Leak_nodes[i].Creation_date)
		fmt.Println("Last Updated:", data.Leak_nodes[i].Update_date)

		httpposturl := os.Getenv("webhook_key")

		fmt.Println("HTTP JSON POST URL:", httpposturl)

		var jsonData = []byte(`{
			"username": "webhooktest",
			"content": "123"
		}`)
		request, error := http.NewRequest("POST", httpposturl, bytes.NewBuffer(jsonData))
		if error != nil {
			fmt.Println("request error")
		}
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")

		fmt.Println(request)

		client := &http.Client{}
		response, error := client.Do(request)
		if error != nil {
			panic(error)
		}
		defer response.Body.Close()

		fmt.Println("response Status:", response.Status)
		//fmt.Println("response Headers:", response.Header)
		body, _ := io.ReadAll(response.Body)
		fmt.Println("response Body:", string(body))

	}

}
