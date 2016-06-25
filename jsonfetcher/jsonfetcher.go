package jsonfetcher

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type FetchedHost struct {
	RunningTask  string `bson:"running_task" json:"running_task"`
	InstanceType string `bson:"instance_type" json:"instance_type,omitempty"`
	Provider     string `bson:"host_type" json:"host_type"`
	Status       string `bson:"status" json:"status"`
}

func FetchPage() []FetchedHost {

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
		result[i] = hostStruct.FetchedHost
	}
	return result
}
