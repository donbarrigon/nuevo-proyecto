package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewHttpServer(port string, routes *Routes) *http.Server {
	timeout := time.Duration(Env.SERVER_TIMEOUT) * time.Second

	router := &Router{}
	router.Make(routes)
	Routers[port] = router

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      router.HandlerFunction(),
		ReadTimeout:  timeout / 2,
		WriteTimeout: timeout / 2,
		IdleTimeout:  timeout,
	}

	go func() {
		startMessage()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			PrintError("Could not start server: :error", Entry{"error", err.Error()})
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
	PrintInfo("Iniciando apagado controlado del servidor...")

	// Crea un contexto con timeout para el shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//se cierra el servidor HTTP para que no acepte nuevas conexiones
	if err := server.Shutdown(ctx); err != nil {
		PrintWarning("Servidor forzado a cerrar: :err", Entry{"err", err.Error()})
	} else {
		PrintInfo("Servidor HTTP detenido correctamente")
	}

	// se cierra la conexion con mono db
	if err := CloseMongoDB(); err != nil {
		PrintWarning("Error al cerrar la conexión a MongoDB: :err", Entry{"err", err.Error()})
	} else {
		PrintInfo("Conexión a MongoDB cerrada correctamente")
	}

	PrintInfo("Apagado controlado completado")
}

func startMessage() {
	PrintInfo(`
   ____   ___  ____  ____  ___  ___  _   _ ____   ___
  / ___| / _ \|  _ \|  _ \|_ _|| __|| \ | |  _ \ / _ \
 | |    | | | | |_) | |_) || ||||__ |  \| | | | | | | |
 | |___ | |_| |  _ <|  _ < | ||||__ | |\  | |_| | |_| |
  \____(_)___/|_| \_\_| \_\___||___||_| \_|____/ \___/

 🚀 Servidor corriendo en http://localhost:{port} 
 🌱 Entorno: DESARROLLO
	`, Entry{"port", Env.SERVER_PORT})
}
