package main

import (
	"log"

	com "github.com/donbarrigon/nuevo-proyecto/internal/common"
	"github.com/donbarrigon/nuevo-proyecto/internal/server"
)

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Iniciando el servidor")

	com.LoadConfig()
	com.InitMongoDB()

	httpServer := server.NewHttpServer(com.SERVER_PORT)
	server.HttpServerGracefulShutdown(httpServer)
}
