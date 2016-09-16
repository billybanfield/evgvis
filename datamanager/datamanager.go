package datamanager

import (
	"log"
	"os"
	"sync"

	fetcher "github.com/billybanfield/evgvis/jsonfetcher"
	"gopkg.in/mgo.v2"
)

const (
	hostsCol = "hosts_state"
)

type sessionManager struct {
	session *mgo.Session
	lock    *sync.RWMutex
}

var globalSession sessionManager

func (s *sessionManager) bulkInsert(inserts []interface{}) {
	session, err := s.getSession()
	lock := s.getLock()
	if err != nil {
		log.Fatalf("error fetching session %v\n", err)
	}

	dbName := os.Getenv("DB_NAME")
	collection := session.DB(dbName).C(hostsCol)
	bulk := collection.Bulk()
	bulk.Insert(inserts...)

	lock.Lock()
	_, err = bulk.Run()
	lock.Unlock()
	if err != nil {
		log.Printf("error updating state: %v\n", err)
	}
}

func (s *sessionManager) removeAll() {
	session, err := s.getSession()
	lock := s.getLock()

	if err != nil {
		log.Fatalf("error fetching session %v\n", err)
	}

	dbName := os.Getenv("DB_NAME")

	lock.Lock()
	session.DB(dbName).C(hostsCol).RemoveAll(nil)
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
	globalSession.removeAll()
	globalSession.bulkInsert(hostsAsInterface)

	log.Println("State updated")
}

func FetchState() []fetcher.FetchedHost {
	log.Print("Fetching State")
	return globalSession.FetchHosts()
}
