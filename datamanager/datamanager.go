package datamanager

import (
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	fetcher "github.com/billybanfield/evgvis/jsonfetcher"
	"gopkg.in/mgo.v2"
)

type EvergreenState struct {
	RunningHosts []fetcher.FetchedHost `json:"running_hosts"`
	ApiStatus    string                `json:"api_status"`
	UiStatus     string                `json:"ui_status"`
}

const (
	ServiceReachable          = "reachable"
	ServiceUnreachable        = "unreachable"
	ServiceReachableWithError = "reachable_with_error"
)

type ServiceStatus struct {
	Id     string `bson:"_id"`
	Status string `bson:"status"`
}

type TimeCounter struct {
	Count     int       `bson:"count"`
	TimeStamp time.Time `bson:"time_stamp"`
}

const (
	hostsCol  = "hosts_state"
	timeCol   = "hosts_over_time"
	statusCol = "status"
)

type sessionManager struct {
	session *mgo.Session
	lock    *sync.RWMutex
}

var globalSession sessionManager

func (s *sessionManager) bulkInsert(inserts []interface{}, collectionName string) {
	session, err := s.getSession()
	if err != nil {
		log.Fatalf("error fetching session %v\n", err)
	}

	dbName := os.Getenv("DB_NAME")
	collection := session.DB(dbName).C(collectionName)
	bulk := collection.Bulk()
	bulk.Insert(inserts...)

	_, err = bulk.Run()
	if err != nil {
		log.Printf("error updating state: %v\n", err)
	}
}

func (s *sessionManager) singleInsert(document interface{}, collectionName string) {
	session, err := s.getSession()
	if err != nil {
		log.Fatalf("error fetching session %v\n", err)
	}

	dbName := os.Getenv("DB_NAME")
	collection := session.DB(dbName).C(collectionName)

	err = collection.Insert(document)
	if err != nil {
		log.Printf("error updating state: %v\n", err)
	}
}

func (s *sessionManager) removeAll(collection string) {
	session, err := s.getSession()

	if err != nil {
		log.Fatalf("error fetching session %v\n", err)
	}
	dbName := os.Getenv("DB_NAME")
	session.DB(dbName).C(collection).RemoveAll(nil)
}
func (s *sessionManager) FetchHosts() []fetcher.FetchedHost {
	session, err := s.getSession()
	lock := s.getLock()
	if err != nil {
		log.Fatalf("error fetching session %v\n", err)
	}
	result := &[]fetcher.FetchedHost{}
	dbName := os.Getenv("DB_NAME")
	lock.Lock()
	session.DB(dbName).C(hostsCol).Find(nil).All(result)
	lock.Unlock()
	return *result
}

func (s *sessionManager) FetchServiceStatus(service string) string {
	session, err := s.getSession()
	lock := s.getLock()
	if err != nil {
		log.Fatalf("error fetching session %v\n", err)
	}
	result := &ServiceStatus{}
	dbName := os.Getenv("DB_NAME")
	serviceFinder := &struct {
		Id string `bson:"_id"`
	}{service}

	lock.Lock()
	session.DB(dbName).C(statusCol).Find(serviceFinder).One(result)
	lock.Unlock()
	return result.Status
}

func (s *sessionManager) getSession() (*mgo.Session, error) {
	uri := os.Getenv("MONGODB_URI")
	if s.session == nil {
		session, err := mgo.Dial(uri)
		if err != nil {
			return nil, err
		}
		s.session = session
	}
	return s.session, nil
}
func (s *sessionManager) getLock() *sync.RWMutex {
	if s.lock == nil {
		s.lock = &sync.RWMutex{}
	}
	return s.lock
}

func UpdateState() {
	log.Print("Updating state")

	hosts := fetcher.FetchHosts()
	hostsAsInterface := make([]interface{}, len(hosts))
	for i := range hosts {
		hostsAsInterface[i] = &hosts[i]
	}
	lock := globalSession.getLock()

	lock.Lock()
	globalSession.removeAll(hostsCol)
	globalSession.bulkInsert(hostsAsInterface, hostsCol)

	timeCounter := &TimeCounter{
		TimeStamp: time.Now(),
		Count:     len(hosts),
	}

	globalSession.singleInsert(timeCounter, timeCol)
	if time.Now().Hour() == 0 {
		globalSession.removeAll(hostsCol)
	}

	apiUrl := os.Getenv("API_URL")
	uiUrl := os.Getenv("UI_URL")

	apiStatus := &ServiceStatus{
		Id: "api",
	}
	uiStatus := &ServiceStatus{
		Id: "ui",
	}
	globalSession.removeAll(statusCol)
	status, err := GetServiceStatus(apiUrl)
	if err != nil {
		panic(err)
	}
	apiStatus.Status = status
	globalSession.singleInsert(apiStatus, statusCol)

	status, err = GetServiceStatus(uiUrl)
	if err != nil {
		panic(err)
	}
	uiStatus.Status = status
	globalSession.singleInsert(uiStatus, statusCol)
	lock.Unlock()

	log.Println("State updated")

}

func GetServiceStatus(url string) (string, error) {
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		if err, ok := err.(net.Error); ok && err.Timeout() {
			return ServiceUnreachable, nil
		}
		return "", err
	}

	if resp.StatusCode != 200 {
		return ServiceReachableWithError, nil
	}

	return ServiceReachable, nil
}

func FetchState() EvergreenState {
	log.Println("Fetching State")

	return EvergreenState{
		RunningHosts: globalSession.FetchHosts(),
		ApiStatus:    globalSession.FetchServiceStatus("api"),
		UiStatus:     globalSession.FetchServiceStatus("ui"),
	}
}
