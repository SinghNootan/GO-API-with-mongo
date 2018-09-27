package main 

import (
	"log"
	"gopkg.in/mgo.v2"
	"fmt"
	"labix.org/v2/mgo/bson"
	"time"
)

type Output struct {
	OputputData Response `json:"response"`
}
type Response struct {
	Status int `json:"status"`
	Data QueryData `json:"data"`
	Message string 	`json:"message"`
}
type QueryData struct {
	Id bson.ObjectId	`json:"id" bson:"_id"`
	First_name  string	`json:"first_name" bson:"first_name"`
	Last_name   string	`json:"last_name" bson:"last_name"`
	Email 		string	`json:"email" bson:"email"`
	Mobile      int 	`json:"mobile" bson:"mobile"`
	Pass        string 	`json:"pass" bson:"pass"`
}

type DbCredentials struct {
	Env string `json:"env"`
	Hosts string `json:"hosts"`
	DatabaseName string `json:"databaseName"`
	Username string `json:"username"`
	Password string `json:"password"`
	Collections string `json:"collections"`
	Ip string `json:"ip"`
	Port string `json:"port"`
}

var DbCred *DbCredentials
var db *mgo.Database
func Connection(){
	fmt.Println("Database connection.....")
	readConfig()
	if DbCred.Env == "dev" {
		session, err := mgo.Dial(DbCred.Hosts)
		if err != nil {
			log.Fatalf("Erorr while connecting mongodb on dev %v", err)
		}
		db = session.DB(DbCred.DatabaseName)
	} else {
		info := &mgo.DialInfo{
	        Addrs:    []string{DbCred.Hosts},
	        Timeout:  60 * time.Second,
	        Database: DbCred.DatabaseName,
	        Username: DbCred.Username,
	        Password: DbCred.Password,
	    }
	    session, err := mgo.DialWithInfo(info)
	    if err != nil {
	        log.Fatalf("Erorr while connecting mongodb on prod %v", err)
	    }
	    db = session.DB(DbCred.DatabaseName)
	}
	
}