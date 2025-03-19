package app

var MONGO_URI string
var DB_NAME string
var PORT string

func LoadConfig() {

	MONGO_URI = "mongodb://localhost:27017"
	DB_NAME = "sample_mflix"
	PORT = "8080"
}
