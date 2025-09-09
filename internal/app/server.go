package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func ServerStart(port string, routes *Routes) {
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

	PrintInfo(`🚀 Server running on :app_url 
  ____   ___  ____  ____  ___  ___  _   _ ____   ___
 / ___| / _ \|  _ \|  _ \|_ _|| __|| \ | |  _ \ / _ \
| |    | | | | |_) | |_) || ||||__ |  \| | | | | | | |
| |___ | |_| |  _ <|  _ < | ||||__ | |\  | |_| | |_| |
 \____(_)___/|_| \_\_| \_\___||___||_| \_|____/ \___/
`, Entry{"app_url", Env.APP_URL})

	// funciona en dev pero en produccion es feo.
	// espera la señal en segundo plano, el bun run dev lo reinicia pero el main se termina y no salen los mensajes
	go HttpServerGracefulShutdown(server)
	time.Sleep(100 * time.Millisecond) // para que salga el mensaje de corriendo.

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		PrintError("🔴💥 Could not start server: :error", Entry{"error", err.Error()})
	}

	// funciona mejor en produccion en dev no.
	// server en segundo plano espera la señal de cierre, salen los mensajes pero el bun run dev no lo reinicia
	// go func() {
	// 	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		PrintError("🔴💥 Could not start server: :error", Entry{"error", err.Error()})
	// 	}
	// }()
	// HttpServerGracefulShutdown(server)

}

// maneja el apagado graceful del servidor
func HttpServerGracefulShutdown(server *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Espera por la señal de terminación
	<-sigChan
	PrintInfo("⏻ Initiating controlled server shutdown...")

	// se cierra la conexion con mono db
	if err := CloseMongoDB(); err != nil {
		PrintWarning("🔴💥 Error closing connection to MongoDB :err", Entry{"err", err.Error()})
	} else {
		PrintInfo("🔌 Connection to MongoDB successfully closed")
	}

	// Crea un contexto con timeout para el shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//se cierra el servidor HTTP para que no acepte nuevas conexiones
	if err := server.Shutdown(ctx); err != nil {
		PrintWarning("⏻ Server forced to close: :err", Entry{"err", err.Error()})
	} else {
		PrintInfo("⏻ HTTP server stopped successfully")
	}

	PrintInfo("💀 Apagado controlado completado")
}
