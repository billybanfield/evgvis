package jsonfetcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type APIHost struct {
	Id          *string    `json:"host_id"`
	Distro      distroInfo `json:"distro"`
	Provisioned bool       `json:"provisioned"`
	Type        *string    `json:"host_type"`
	Status      *string    `json:"status"`
	RunningTask taskInfo   `json:"running_task"`
}

type distroInfo struct {
	Id       *string `json:"distro_id"`
	Provider *string `json:"provider"`
}

type taskInfo struct {
	Id           *string `json:"task_id"`
	Name         *string `json:"name"`
	DispatchTime *string `json:"dispatch_time"`
	VersionId    *string `json:"version_id"`
	BuildId      *string `json:"build_id"`
}

type APIError struct {
	Message    string `json:"error"`
	StatusCode int    `json:"status"`
}

func FetchHosts() []APIHost {

	log.Println("fetching hosts")
	uiUrl := os.Getenv("UI_URL")
	runningHostsUrl := fmt.Sprintf("%vrest/v2/hosts?status=running&limit=100000000", uiUrl)

	out1, err := hostRequestHelper("GET", runningHostsUrl)
	if err != nil {
		log.Printf("error fetching host data: ", err)
	}
	decoHostsUrl := fmt.Sprintf("%vrest/v2/hosts?status=decommissioned&limit=100000000", uiUrl)
	out2, err := hostRequestHelper("GET", decoHostsUrl)
	if err != nil {
		log.Printf("error fetching host data: ", err)
	}
	startingHostsUrl := fmt.Sprintf("%vrest/v2/hosts?status=starting&limit=100000000", uiUrl)
	out3, err := hostRequestHelper("GET", startingHostsUrl)
	if err != nil {
		log.Printf("error fetching host data: ", err)
	}
	unreachableHostsUrl := fmt.Sprintf("%vrest/v2/hosts?status=unreachable&limit=100000000", uiUrl)
	out4, err := hostRequestHelper("GET", unreachableHostsUrl)
	if err != nil {
		log.Printf("error fetching host data: ", err)
	}
	quarantinedHostsUrl := fmt.Sprintf("%vrest/v2/hosts?status=quarantined&limit=100000000", uiUrl)
	out5, err := hostRequestHelper("GET", quarantinedHostsUrl)
	if err != nil {
		log.Printf("error fetching host data: ", err)
	}
	total := []APIHost{}
	total = append(total, out1...)
	total = append(total, out2...)
	total = append(total, out3...)
	total = append(total, out4...)
	total = append(total, out5...)

	return total
}

func hostRequestHelper(method, url string) ([]APIHost, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Api-Key", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Auth-Username", os.Getenv("AUTH_USER"))

	resp, err := client.Do(req)
	if err != nil {
		return []APIHost{}, err
	}
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return []APIHost{}, nil
		}
		bytes, _ := ioutil.ReadAll(resp.Body)
		out := APIError{}
		err = json.Unmarshal(bytes, &out)
		if err != nil {
			log.Fatalf("error unmarshaling data: %v", err)
			return []APIHost{}, err
		}
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	out := []APIHost{}
	err = json.Unmarshal(bytes, &out)
	if err != nil {
		log.Fatalf("error unmarshaling data: %v", err)
		return []APIHost{}, err
	}
	return out, nil
}
