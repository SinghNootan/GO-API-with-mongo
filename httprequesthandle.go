package main

import (
	"log"
	"net/http"
	"fmt"
	"labix.org/v2/mgo/bson"
	"encoding/json"
	//"strconv"
	"os"
	"io/ioutil"
)

func readConfig() {
	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatalf("Erro while opening config file to read %v",err)
	} else {
		byteslice, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatalf("Error while reading file to read conf file: %v", err)
		} else {
			err = json.Unmarshal(byteslice, &DbCred)
			if err != nil {
				log.Fatalf("error while unmarshalling conf file: %v", err)
			}
		}
	}
	defer file.Close()
}

func ResponseWithError(w http.ResponseWriter, msg string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var qd QueryData
	if msg != "" {
		var FinalData Output
		FinalData.OputputData.Status = 500
		FinalData.OputputData.Message = msg
		FinalData.OputputData.Data = qd
		res, err := json.Marshal(FinalData)
		if err != nil {
			log.Fatalf("Error while marshalling response error data %v", err)
		}
		fmt.Fprintf(w, string(res))
	}
}

func ResponseWithJson(w http.ResponseWriter, data QueryData) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	var FinalData Output
	FinalData.OputputData.Status = 200
	FinalData.OputputData.Message = "Success"
	FinalData.OputputData.Data = data
	res, err := json.Marshal(FinalData)
	if err != nil {
		log.Fatalf("Error while marshalling response error data %v", err)
	}
	fmt.Fprintf(w, string(res))
}


func loginUser(w http.ResponseWriter, r *http.Request) {
	//var response string
	email, ok := r.URL.Query()["email"]
	defer r.Body.Close()
	if !ok || len(email[0]) < 1 {
		ResponseWithError(w, "email is missing")
	} else {
		pass, ok := r.URL.Query()["pass"]
		if !ok || len(pass[0]) < 1 {
			ResponseWithError(w, "pass is missing")
		} else {
			emailId := string(email[0])
			password := string(pass[0])
			c := db.C(DbCred.Collections)
			var result QueryData
			err := c.Find(bson.M{"email":emailId,"pass":password}).One(&result)
			if err != nil {
				log.Printf("Erorr while fetching data %v", err)
				ResponseWithError(w, "Erorr while fetching data")
			} else {
				ResponseWithJson(w, result)
			}
		}
	}
}
func getAllList(w http.ResponseWriter, r *http.Request) {
	c := db.C(DbCred.Collections)
	var result []QueryData
	//err = c.Find(bson.M{"first_name":"nootan"}).All(&result)
	err := c.Find(nil).All(&result)
	if err != nil {
		ResponseWithError(w, "Erorr while fetching data")
	} else {
		response, err := json.Marshal(result)
		if err != nil {
			//log.Fatalf("Error while marshalling result of query %v", err)
			ResponseWithError(w, "Erorr while fetching data")
		} else {
			fmt.Fprintf(w, string(response))
		}
	}
} 
func registerUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var newdata QueryData
	if err := json.NewDecoder(r.Body).Decode(&newdata); err != nil {
		ResponseWithError(w, "invalid request payload")
	} else {
		c := db.C(DbCred.Collections)
		var result QueryData
		err := c.Find(bson.M{"email":newdata.Email,"pass":newdata.Pass}).One(&result)
		if err != nil {
			newdata.Id = bson.NewObjectId()
			fmt.Println(newdata)
			err := c.Insert(&newdata)
			if err != nil {
				ResponseWithError(w, "Erorr while inserting data")
			} else {
				ResponseWithJson(w, newdata)
			}
		} else {
			fmt.Println(result)
			if result != (QueryData{}) {
				ResponseWithError(w, "User already exists")
			}
		}
	}
}

func main() {
	readConfig()
	h := http.NewServeMux()
	Connection()
	h.HandleFunc("/", getAllList)
	h.HandleFunc("/login", loginUser)
	h.HandleFunc("/create", registerUser)
	//h.HandleFunc("/symbols", getSymbols)
	//h.HandleFunc("/history", getAlldata)
	hostWithPort := DbCred.Ip+":"+DbCred.Port
	err := http.ListenAndServe(hostWithPort, h)
	fmt.Println(err)
}
