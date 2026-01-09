package controller

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type AppServer struct {
	http.Server
}

func Start(srv *http.Server) {
	log.Printf("sever starting at address %s", srv.Addr)

	// シグナル監視用 Context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// サーバ起動
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("server error: %v", err)
		}
	}()

	// シグナル待ち
	<-ctx.Done()
	log.Println("shutdown signal received")

	// graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	log.Println("server gracefully stopped")
}
