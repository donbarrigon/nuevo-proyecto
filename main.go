package main

import (
	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/routes"
)

func main() {
	app.LoadEnv()
	app.InitMongoDB()
	app.ServerStart(app.Env.SERVER_PORT, routes.GetAll())
}
