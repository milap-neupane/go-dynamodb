package main

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type Worker struct {
	Id        string `json:id`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Address   string `json:"address,omitempty"`
}

func createWorker(worker *Worker) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String("us-east-1")},
		SharedConfigState: session.SharedConfigEnable,
		// Profile: "",
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(worker)

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Workers"),
	}

	resp, err := svc.PutItem(input)
	log.Println(resp)
	if err != nil {
		log.Println("Got error calling PutItem:")
		log.Println(err.Error())
		os.Exit(1)
	}

	log.Println("Successfully added Worker")
}

func getWorker(worker *Worker) *Worker {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Region: aws.String("us-east-1")},
		SharedConfigState: session.SharedConfigEnable,
		Profile:           "team0workermgmt",
	}))

	svc := dynamodb.New(sess)
	log.Println(worker)
	key, err := dynamodbattribute.MarshalMap(worker)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	input := &dynamodb.GetItemInput{
		Key:       key,
		TableName: aws.String("Workers"),
	}

	result, err := svc.GetItem(input)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	movie := Worker{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &movie)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return &movie
}

func main() {
	http.HandleFunc("/worker", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			id := r.URL.Query().Get("Id")
			worker := getWorker(&Worker{Id: string(id)})
			json.NewEncoder(w).Encode(worker)
		case "POST":
			var worker Worker
			if err := json.NewDecoder(r.Body).Decode(&worker); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			uid, _ := exec.Command("uuidgen").Output()
			worker.Id = string(uid)
			createWorker(&worker)
			json.NewEncoder(w).Encode(map[string]string{"id": string(uid)})
		default:
			w.WriteHeader(http.StatusNotImplemented)
			w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
		}
	})
	http.ListenAndServe(":3000", nil)
}
