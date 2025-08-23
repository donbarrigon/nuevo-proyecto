package seed

import "github.com/donbarrigon/nuevo-proyecto/internal/app"

var Seeds = app.Object{}

func Run() {
	// inserte las funciones de seed() carguelas todas que despues el comando run seed ejecuta solo las que no estan cargadas
	add("World", World)
}

func add(name string, fun func()) {
	Seeds.Set(name, fun)
}
