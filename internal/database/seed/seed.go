package seed

import "github.com/donbarrigon/nuevo-proyecto/internal/app"

var Seeds = app.Fields{}

func Run() {
	// inserte las funciones de seed() carguelas todas que despues el comando run seed ejecuta solo las que no estan cargadas
	add("Cities", Cities)
}

func add(name string, fun func()) {
	Seeds = append(Seeds, app.F{Key: name, Value: fun})
}
