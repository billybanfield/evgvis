package jsonfetcher

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type FetchedHost struct {
	Id           string      `bson:"_id" json:"id"`
	RunningTask  string      `bson:"running_task" json:"running_task"`
	InstanceType string      `bson:"instance_type" json:"instance_type"`
	Provider     string      `bson:"host_type" json:"host_type"`
	Status       string      `bson:"status" json:"status"`
	Distro       interface{} `bson:"distro" json:distro"`
}

func FetchHosts() []FetchedHost {

	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://mci-motu.10gen.cc:9090/hosts", nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error getting data: %v", err)
	}
	hostsRegexp := regexp.MustCompile("window.hosts =.*")

	bytes, _ := ioutil.ReadAll(resp.Body)
	found := hostsRegexp.Find(bytes)
	marshalled := found[15 : len(found)-1]
	out := &struct {
		Hosts []struct {
			FetchedHost `json:"Host"`
		} `json:"Hosts"`
	}{}

	err = json.Unmarshal(marshalled, out)
	if err != nil {
		log.Fatalf("error unmarshaling data: %v", err)
	}

	result := make([]FetchedHost, len(out.Hosts))
	for i, hostStruct := range out.Hosts {
		if distro, ok := hostStruct.FetchedHost.Distro.(map[string]interface{}); ok {
			hostStruct.FetchedHost.Distro, ok = distro["_id"].(string)
			if !ok {
				panic("distro id not string")
			}
		}
		result[i] = hostStruct.FetchedHost
	}
	return result
}
