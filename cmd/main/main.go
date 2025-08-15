package main

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
	"github.com/donbarrigon/nuevo-proyecto/internal/http/routes"
)

func main() {

	app.LoadEnv()

	db.InitMongoDB()

	httpServer := app.NewHttpServer(app.Env.SERVER_PORT, routes.GetApi())

	app.HttpServerGracefulShutdown(httpServer)
	db.CloseMongoConnection()
}
