package main

import (
	"clientgo/config"
	"clientgo/db"
	"clientgo/routes"
	"clientgo/tcpclient"
	"clientgo/tcpserver"
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Load configuration
	config.GetConfigurations()

	// Connect to Redis
	db.ConnectRedisClient()
	defer db.CloseRedis()

	
	runClient := flag.Bool("client", false, "Start TCP client")
	runServers := flag.Bool("tcp-servers", false, "Start 10 TCP servers")
	runAll := flag.Bool("all", false, "Start both TCP servers and client")
	flag.Parse()

	switch {
	case *runAll:
		log.Println("Starting all TCP servers and client...")
		go tcpserver.StartMultipleTCPServers()
		time.Sleep(2 * time.Second) 
		tcpclient.StartClient()
		return

	case *runServers:
		tcpserver.StartMultipleTCPServers()
		return

	case *runClient:
		tcpclient.StartClient()
		return

	default:
		// Start HTTP API server
		if err := runHTTPServer(); err != nil {
			log.Fatalln(err)
		}
	}

	// Start REST HTTP server
	if err := runHTTPServer(); err != nil {
		log.Fatalln(err)
	}

}

func runHTTPServer() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	router := gin.Default()
	routes.RegisterRouter(router)

	srv := &http.Server{
		Addr:         ":" + config.AppConfig.Server.Port,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}

	errorChan := make(chan error, 1)
	go func() {
		errorChan <- srv.ListenAndServe()
	}()

	log.Printf("HTTP Server running on http://localhost:%s/ping\n", config.AppConfig.Server.Port)

	select {
	case err := <-errorChan:
		return err
	case <-ctx.Done():
		log.Println("Server interrupt received, shutting down...")
		stop()
		return srv.Shutdown(context.Background())
	}
}
