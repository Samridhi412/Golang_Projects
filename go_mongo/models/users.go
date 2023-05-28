package models

import (
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id bson.ObjectId `json:"id" bson:"_id"` //id is stored as an underscore in bson
	Name string `json:"name" bson:"name"`
	Gender string `json:"gender" bson:"gender"`
	Age int `json:"age" bson:"age"`
}