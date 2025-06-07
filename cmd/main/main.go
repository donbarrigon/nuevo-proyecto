package main

import (
	"log"

	"github.com/donbarrigon/nuevo-proyecto/internal/config"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	"github.com/donbarrigon/nuevo-proyecto/internal/server"
)

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Iniciando el servidor")

	if err := config.Load(); err != nil {
		log.Println(err.Error())
		return
	}
	db.InitMongoDB()

	httpServer := server.NewHttpServer(config.Env.SERVER_PORT)
	server.HttpServerGracefulShutdown(httpServer)
}
