package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/database/db"
)

func NewHttpServer(port string) *http.Server {
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      NewRouter(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", port)
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
