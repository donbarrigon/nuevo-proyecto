package main

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	"github.com/donbarrigon/nuevo-proyecto/internal/server"
)

func main() {

	app.LoadEnv()

	db.InitMongoDB()

	httpServer := server.NewHttpServer(app.Env.SERVER_PORT)
	server.HttpServerGracefulShutdown(httpServer)
}
