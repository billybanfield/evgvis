package datamanager

import (
	"log"
	"os"
	"sync"
	"time"

	fetcher "github.com/billybanfield/evgvis/jsonfetcher"
	"gopkg.in/mgo.v2"
)

type TimeCounter struct {
	Count     int       `bson:"count"`
	TimeStamp time.Time `bson:"time_stamp"`
}

const (
	hostsCol = "hosts_state"
	timeCol  = "hosts_over_time"
)

type sessionManager struct {
	session *mgo.Session
	lock    *sync.RWMutex
}

var globalSession sessionManager

func (s *sessionManager) bulkInsert(inserts []interface{}, collectionName string) {
	session, err := s.getSession()
	lock := s.getLock()
	if err != nil {
		log.Fatalf("error fetching session %v\n", err)
	}

	dbName := os.Getenv("DB_NAME")
	collection := session.DB(dbName).C(collectionName)
	bulk := collection.Bulk()
	bulk.Insert(inserts...)

	lock.Lock()
	_, err = bulk.Run()
	lock.Unlock()
	if err != nil {
		log.Printf("error updating state: %v\n", err)
	}
}

func (s *sessionManager) singleInsert(document interface{}, collectionName string) {
	session, err := s.getSession()
	lock := s.getLock()
	if err != nil {
		log.Fatalf("error fetching session %v\n", err)
	}

	dbName := os.Getenv("DB_NAME")
	collection := session.DB(dbName).C(collectionName)

	lock.Lock()
	err = collection.Insert(document)
	lock.Unlock()
	if err != nil {
		log.Printf("error updating state: %v\n", err)
	}
}

func (s *sessionManager) removeAll(collection string) {
	session, err := s.getSession()
	lock := s.getLock()

	if err != nil {
		log.Fatalf("error fetching session %v\n", err)
	}

	dbName := os.Getenv("DB_NAME")

	lock.Lock()
	session.DB(dbName).C(collection).RemoveAll(nil)
	lock.Unlock()
}
func (s *sessionManager) FetchHosts() []fetcher.FetchedHost {
	session, err := s.getSession()
	lock := s.getLock()
	if err != nil {
		log.Fatalf("error fetching session %v\n", err)
	}
	result := &[]fetcher.FetchedHost{}
	dbName := os.Getenv("DB_NAME")
	lock.RLock()
	session.DB(dbName).C(hostsCol).Find(nil).All(result)
	lock.RUnlock()
	return *result
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
	log.Println("State updated")

}

func FetchState() []fetcher.FetchedHost {
	log.Println("Fetching State")
	return globalSession.FetchHosts()
}
