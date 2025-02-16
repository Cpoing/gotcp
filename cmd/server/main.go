package main

import (
  "net"
  "fmt"
  "log"
  "os"
  "os/signal"
  "syscall"

  "github.com/Cpoing/gotcp/pkg/config"
  "github.com/Cpoing/gotcp/pkg/chat"
  "github.com/Cpoing/gotcp/pkg/connection"
)

func main() {

  cfg, err := config.LoadConfig("config.json")
  if err != nil {
    log.Fatalf("Failed to load configuration: %v", err)
  }

  log.Println("Starting TCP Chat Server...")

  chatManager := chat.NewManager()
  address := fmt.Sprintf(":%d", cfg.Port)

  listener, err := net.Listen("tcp", address)
  if err != nil {
    fmt.Println("Error: ", err)
    return
  }
  defer listener.Close()

  log.Printf("Server is listening on port %s", address)

  quit := make(chan os.Signal, 1)
  signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
  go func() {
    <-quit
    log.Println("Shutting down server...")
    listener.Close()
    os.Exit(0)
  }()
  
  for {
    conn, err := listener.Accept()
    if err != nil {
      log.Printf("Error accepting connection: %v", err)
      continue
    }

    go connection.HandleConnection(conn, chatManager)
  }
}
