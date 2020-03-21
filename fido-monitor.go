package main

import (
	"context"
	"encoding/json"
	"fido-monitor/structs"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var homeHTML, _ = template.ParseFiles("./views/index.html")
var tz, _ = time.LoadLocation("Asia/Jakarta")

func main() {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		fmt.Println("here")
		panic(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	fmt.Println("connecting...")
	err = client.Connect(ctx)
	fmt.Println("connected...")
	fmt.Println(err)

	if err != nil {
		panic(err)
	}

	collection = client.Database("fido_meta").Collection("link")

	http.HandleFunc("/", Home)
	http.HandleFunc("/api/link-by-date", GetLinkByDate)
	http.HandleFunc("/api/link-latest", GetLatestLink)
	fmt.Println("Starting Server on :8080")
	http.ListenAndServe(":8080", nil)
}

// Home main function that serve home page
func Home(w http.ResponseWriter, r *http.Request) {
	homeHTML.Execute(w, nil)
}

// GetLinkByDate main controller of this API
func GetLinkByDate(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	query := r.URL.Query()

	layoutFormat := "02-01-2006"
	startDate, _ := time.ParseInLocation(layoutFormat, query.Get("start"), tz)
	endDate, _ := time.ParseInLocation(layoutFormat, query.Get("end"), tz)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := collection.Find(ctx, bson.M{"_id": bson.M{"$gte": ObjID(startDate), "$lte": ObjID(endDate)}})
	if err != nil {
		// https://stackoverflow.com/a/40096757/12985309
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened! " + err.Error()))
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)

	var result []structs.Link

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	for cur.Next(ctx) {
		var doc structs.Link
		err := cur.Decode(&doc)
		if err != nil {
			log.Fatal(err)
		}
		// do something with result....
		result = append(result, doc)
	}

	m, _ := json.Marshal(result)
	fmt.Fprint(w, string(m))
}

// GetLatestLink only get latest link with limit
func GetLatestLink(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var opt options.FindOptions
	opt.SetBatchSize(15)
	opt.SetLimit(15)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cur, err := collection.Find(ctx, bson.M{}, &opt)
	if err != nil {
		// https://stackoverflow.com/a/40096757/12985309
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened! " + err.Error()))
		return
	}

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	defer cur.Close(ctx)

	var result []structs.Link

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	for cur.Next(ctx) {
		var doc structs.Link
		err := cur.Decode(&doc)
		if err != nil {
			log.Fatal(err)
		}
		// do something with result....
		result = append(result, doc)
	}

	m, _ := json.Marshal(result)
	fmt.Fprint(w, string(m))
}

// ObjID simplify function name of NewObjectIDFromTimestamp
func ObjID(date time.Time) primitive.ObjectID {
	return primitive.NewObjectIDFromTimestamp(date)
}
