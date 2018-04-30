package dao

import (
	"github.com/elastic/beats/libbeat/logp"

	"github.com/jcsuscriptor/elmahbeat/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ElmahDAO struct {
	Server   string
	Database string
	Collection string
}

var db *mgo.Database


// Establish a connection to database
func (m *ElmahDAO) Connect() {
 
	logp.Debug("ElmahDAO","Connecting  %s. Database %s", m.Server, m.Database)
	session, err := mgo.Dial(m.Server)
	if err != nil {
		logp.Err(err.Error())
	}
	db = session.DB(m.Database)
	logp.Info("Connected  %s. Database %s", m.Server, m.Database)
}

// Find list of error elmah
func (m *ElmahDAO) FindAll() ([]models.ElmahError, error) {
	
	logp.Debug("ElmahDAO","FindAll %s",m.Collection)

	var data []models.ElmahError

	err := db.C(m.Collection).Find(bson.M{}).All(&data)
	
	if err != nil {
		logp.Err(err.Error())
	}
	return data, err
}
