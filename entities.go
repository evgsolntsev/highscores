package main

type Score struct {
	ID    string  `json:"_id" bson:"_id"`
	Name  string  `json:"name" bson:"name"`
	Score float64 `json:"score" bson:"score"`
}
