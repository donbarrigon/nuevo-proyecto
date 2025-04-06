package controller

import (
	"fmt"
	"net/http"
	"time"
)

type ControllerPuntero struct{}

func (p *ControllerPuntero) Show(ctx *Context) {
	ctx.Writer.Write([]byte("."))
}

type ControllerVacio struct{}

func (ControllerVacio) Show(ctx *Context) {
	ctx.Writer.Write([]byte("."))
}

func Prueba(w http.ResponseWriter, r *http.Request) {
	ctx := NewContext(w, r)
	const iterations = 1_000_000

	// Caso puntero: se crea una nueva instancia cada vez
	benchmark("ControllerPuntero (nueva instancia cada vez)", func() {
		(&ControllerPuntero{}).Show(ctx)
	}, iterations)

	// Caso valor: se crea una nueva instancia cada vez
	benchmark("ControllerVacio (nueva instancia cada vez)", func() {
		(ControllerVacio{}).Show(ctx)
	}, iterations)

	// Opcional: evitar nueva instancia cada vez
	p := &ControllerPuntero{}
	v := ControllerVacio{}

	benchmark("ControllerPuntero (instancia reutilizada)", func() {
		p.Show(ctx)
	}, iterations)

	benchmark("ControllerVacio (instancia reutilizada)", func() {
		v.Show(ctx)
	}, iterations)
}

func benchmark(label string, f func(), iterations int) {
	start := time.Now()
	for i := 0; i < iterations; i++ {
		f()
	}
	duration := time.Since(start)
	fmt.Printf("%s: %s\n", label, duration)
}
