package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"social-app/docs"
	"social-app/internal/connector"
	"social-app/internal/routes"
	"social-app/pkg/ws"
)

type Server struct {
	bc      ws.Broadcaster
	rootCtx context.Context
	root    *gin.Engine
	srv     *http.Server
	rc      *connector.RedisConnector
}

func NewServer(ctx context.Context, router routes.Router, bc ws.Broadcaster, rc *connector.RedisConnector) Server {
	r := gin.New()
	r.MaxMultipartMemory = 8 << 20 // 8 MB
	docs.SwaggerInfo.BasePath = "/"
	router.SetupRoutes(r)
	return Server{
		root: r,
		srv: &http.Server{
			Addr:         ":3222",
			Handler:      r,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 60 * time.Second,
		},
		rootCtx: ctx,
		bc:      bc,
		rc:      rc,
	}
}

func (s *Server) Serve() error {
	ctx, stop := signal.NotifyContext(s.rootCtx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer func() {
		stop()
		fmt.Println(ctx, "🔔 Context cancelled; stopping server...")
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Failed to start server: %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.bc.Start(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				if err := s.rc.Close(); err != nil {
					log.Printf("Error closing Redis client: %v", err)
					return
				}
				log.Println("Redis client closed successfully")
				return
			default:
				err := s.rc.Ping(ctx)
				if err == nil {
					time.Sleep(3 * time.Second)
					continue
				}

				log.Printf("Failed to connect to Redis: %v", err)
				time.Sleep(5 * time.Second)
			}
		}
	}()

	fmt.Println("Server is running on", s.srv.Addr)
	<-ctx.Done()

	fmt.Println(ctx, "🔔 Interrupt signal received; shutting down servers...")

	if err := s.Close(); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	wg.Wait()

	fmt.Println(ctx, "✅ Servers shutdown completed")

	return nil
}

func (s *Server) Close() error {
	if err := s.srv.Shutdown(s.rootCtx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	fmt.Println("Server shutdown gracefully")
	return nil
}
