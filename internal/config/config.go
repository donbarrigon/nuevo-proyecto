package config

var MONGO_URI string
var DB_NAME string
var SERVER_PORT string
var LANG string

func LoadConfig() {

	MONGO_URI = "mongodb+srv://donbarrigon:me5zDK7dfKhmXpxL@cluster0.uu7k0.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	DB_NAME = "sample_mflix"
	SERVER_PORT = "8080"
	LANG = "es"
}
