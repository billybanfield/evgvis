package datamanager

import (
	"log"
	"os"

	"github.com/billybanfield/heroku2/jsonfetcher"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	hostsCol = "hosts_state"
)

type sessionManager struct {
	session *mgo.Session
}

var globalSession sessionManager

func (s *sessionManager) GetSession() (*mgo.Session, error) {
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

func UpdateState() {
	log.Print("Updating state")

	session, err := globalSession.GetSession()
	if err != nil {
		log.Fatal("erroring fetching session %v\n", err)
	}
	dbName := os.Getenv("DB_NAME")
	collection := session.DB(dbName).C(hostsCol)

	collection.RemoveAll(&bson.D{})

	hosts := jsonfetcher.FetchHosts()
	hostsAsInterface := make([]interface{}, len(hosts))
	for i, host := range hosts {
		hostsAsInterface[i] = &host
	}

	bulk := collection.Bulk()
	bulk.Insert(hostsAsInterface...)
	_, err = bulk.Run()
	if err != nil {
		log.Fatal("erroring updating state: %v\n", err)
	}
	log.Println("State updated")

}
