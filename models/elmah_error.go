package models

import "time"
import "gopkg.in/mgo.v2/bson"


type ElmahError struct {
	ID    bson.ObjectId `bson:"_id" json:"id"`
	Application        string        `bson:"ApplicationName" json:"ApplicationName"`
	Host  string        `bson:"host" json:"host"`
	Type  string        `bson:"type" json:"type"`
	Message  string        `bson:"message" json:"message"`
	Source  string        `bson:"source" json:"source"`
	Detail  string        `bson:"detail" json:"detail"`
	User  string        `bson:"user" json:"user"`
	Time    time.Time      `bson:"time" json:"time"`
	StatusCode  int          `bson:"statusCode" json:"statusCode"`
	WebHostHtmlMessage  string        `bson:"webHostHtmlMessage" json:"webHostHtmlMessage"`
	//ServerVariables     bson.M `bson:"ServerVariables,inline"  json:"serverVariables"`

}

 