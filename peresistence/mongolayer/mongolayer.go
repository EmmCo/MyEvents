package mongolayer

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DB     = "myevents"
	USERS  = "users"
	EVENTS = "events"
)

type MongoDBLayer struct {
	session *mgo.Session
}

func NewMongoDBLayer(connection string) (*MongoDBLayer, error) {
	s, err := mgo.Dial(connection)
	if err != nil {
		return nil, err
	}
	return &MongoDBLayer{
		session: s,
	}, nil
}

func (mgolayer *MongoDBLayer) getFreshSession() *mgo.Session {
	return mgolayer.session.Copy()
}

func (mgolayer *MongoDBLayer) FindEvent(id []byte) (Event, error) {
	s := mgolayer.getFreshSession()
	defer s.Close()
	e := Event{}
	err := s.DB(DB).C(EVENTS).FindId(bson.ObjectId(id)).One(&e)
	return e, err
}

func (mgolayer *MongoDBLayer) FindAllAvailableEvents() ([]Event, error) {
	s := mgolayer.getFreshSession()
	defer s.Close()
	es := []Event{}
	err := s.DB(DB).C(EVENTS).Find(nil).All(&es)
	return es, err
}

func (mgolayer *MongoDBLayer) FindEventByName(name string) (Event, error) {
	s := mgolayer.getFreshSession()
	defer s.Close()
	e := Event{}
	err := s.DB(DB).C(EVENTS).Find(bson.M{"name": name}).One(&e)
	return e, err
}

func (mgolayer *MongoDBLayer) AddEvent(e Event) ([]byte, error) {
	s := mgolayer.getFreshSession()
	defer s.Close()
	if !e.ID.Valid() {
		e.ID = bson.NewObjectId()
	}
	if !e.Location.ID.Valid() {
		e.Location.ID = bson.NewObjectId()
	}
	return []byte(e.ID), s.DB(DB).C(EVENTS).Insert(e)
}
