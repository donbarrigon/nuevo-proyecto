package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
)

func NewHttpServer(port string) *http.Server {
	timeout := time.Duration(app.Env.SERVER_TIMEOUT) * time.Second
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      NewRouter(),
		ReadTimeout:  timeout / 2,
		WriteTimeout: timeout / 2,
		IdleTimeout:  timeout,
	}

	go func() {
		startMessage()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %v\n", err)
		}
	}()

	return server
}

// maneja el apagado graceful del servidor
func HttpServerGracefulShutdown(server *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Espera por la señal de terminación
	<-sigChan
	log.Println("Iniciando apagado controlado del servidor...")

	// Crea un contexto con timeout para el shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//se cierra el servidor HTTP para que no acepte nuevas conexiones
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Servidor forzado a cerrar: %v\n", err)
	} else {
		log.Println("Servidor HTTP detenido correctamente")
	}

	// se cierra la conexion con mono db
	if err := db.CloseMongoConnection(); err != nil {
		log.Printf("Error al cerrar la conexión a MongoDB: %v\n", err)
	} else {
		log.Println("Conexión a MongoDB cerrada correctamente")
	}

	log.Println("Apagado controlado completado")
}

func startMessage() {
	log.Println(fmt.Sprintf(`
   ____   ___  ____  ____  ___  ___  _   _ ____   ___
  / ___| / _ \|  _ \|  _ \|_ _|| __|| \ | |  _ \ / _ \
 | |    | | | | |_) | |_) || ||||__ |  \| | | | | | | |
 | |___ | |_| |  _ <|  _ < | ||||__ | |\  | |_| | |_| |
  \____(_)___/|_| \_\_| \_\___||___||_| \_|____/ \___/

 🚀 Servidor corriendo en http://localhost:%v
 🌱 Entorno: DESARROLLO
	`, app.Env.SERVER_PORT))
}
