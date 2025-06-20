package main

import (
	"log"

	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	"github.com/donbarrigon/nuevo-proyecto/internal/server"
	"github.com/donbarrigon/nuevo-proyecto/pkg/system"
)

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Iniciando el servidor")

	if err := system.LoadEnv(); err != nil {
		log.Println(err.Error())
		return
	}
	db.InitMongoDB()

	httpServer := server.NewHttpServer(system.Env.SERVER_PORT)
	server.HttpServerGracefulShutdown(httpServer)
}
