package main

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/http/routes"
)

func main() {

	app.LoadEnv()

	app.InitMongoDB()

	httpServer := app.NewHttpServer(app.Env.SERVER_PORT, routes.GetApi())

	app.HttpServerGracefulShutdown(httpServer)
}
