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
		ip := data.Leak_nodes[i].Ip
		resource_id := data.Leak_nodes[i].Resource_id
		leaks_found := data.Leak_nodes[i].Leak_count
		total_leaks := data.Leak_nodes[i].Leak_event_count
		open_ports := data.Leak_nodes[i].Open_ports

		for x := 0; x < len(data.Leak_nodes[i].Events); x++ {
			e_summary := data.Leak_nodes[i].Events[x].Summary
			e_time := data.Leak_nodes[i].Events[x].Time
			e_event_source := data.Leak_nodes[i].Events[x].Event_source
			e_host_ip := data.Leak_nodes[i].Events[x].Ip
			e_host_address := data.Leak_nodes[i].Events[x].Host
			e_host_ports := data.Leak_nodes[i].Events[x].Port
			e_protocol := data.Leak_nodes[i].Events[x].Protocol
			e_root := data.Leak_nodes[i].Events[x].HTTP.Root
			e_http_url := data.Leak_nodes[i].Events[x].HTTP.URL
			e_organisation := data.Leak_nodes[i].Events[x].Network.Organisation

		}

		first_seen := data.Leak_nodes[i].Creation_date
		last_updated := data.Leak_nodes[i].Update_date

		webhook_url := os.Getenv("webhook_key")

		fmt.Println("HTTP JSON POST URL:", webhook_url)

		values := map[string]string{
			"username": "WebHook Test",
			"content":  ip,
		}

		jsonData, error := json.Marshal(values)
		if error != nil {
			panic(error)
		}

		fmt.Printf("\nThis is the value: %s \n", ip)

		request, error := http.NewRequest("POST", webhook_url, bytes.NewBuffer(jsonData))
		if error != nil {
			fmt.Println("request error")
		}
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")

		client := &http.Client{}
		response, error := client.Do(request)
		if error != nil {
			panic(error)
		}
		defer response.Body.Close()

		fmt.Println("response Status:", response.Status)
		body, _ := io.ReadAll(response.Body)
		fmt.Println("response Body:", string(body))

	}

}
