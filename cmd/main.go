package main

import (
	"log"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/server"
)

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Iniciando el servidor")

	app.LoadConfig()
	app.InitMongoDB()

	httpServer := server.NewHttpServer(app.PORT)
	server.HttpServerGracefulShutdown(httpServer)
}
