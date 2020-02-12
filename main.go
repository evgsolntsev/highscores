package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo"
)

var (
	mySigningKey []byte
	dao          DAO
)

func Ping(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Pong.\n"))
}

func Top(w http.ResponseWriter, req *http.Request) {
	q, ok := req.URL.Query()["q"]
	quantity := 10
	if ok && len(q[0]) > 0 {
		fmt.Sscanf(q[0], "%d", &quantity)
	}

	scores, err := dao.GetTop(context.TODO(), quantity)
	if err != nil {
		log.Printf("Failed to get scores for quantity %v\n", quantity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(scores)
	if err != nil {
		log.Printf("Failed to marshal json on response")
		return
	}

	w.Write(response)

}

func Add(w http.ResponseWriter, req *http.Request) {
	if req.Header["Token"] != nil {
		token, err := jwt.Parse(req.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return mySigningKey, nil
		})

		if err != nil {
			log.Printf("Error adding new score: %v", err.Error())
		}

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			newScore := Score{
				ID:    claims["_id"].(string),
				Name:  claims["name"].(string),
				Score: claims["score"].(float64),
			}
			dao.Insert(context.TODO(), newScore)
			
		}
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Not Authorized.\n"))
	}
}

type Configuration struct {
	Shared   string `json:"shared"`
	MongoURL string `json:"mongo"`
}

const CONFIGNAME = "conf.json"

func init() {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}

	err := decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}
	session, err := mgo.Dial(configuration.MongoURL)
	if err != nil {
		panic(err)
	}

	dao = NewDAO(context.TODO(), session)
	mySigningKey = []byte(configuration.Shared)
}

func main() {
	http.HandleFunc("/ping", Ping)
	http.HandleFunc("/top", Top)
	http.HandleFunc("/add", Add)

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
