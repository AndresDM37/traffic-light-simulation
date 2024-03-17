package main

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "time"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan map[string]string) // broadcast channel

func handleMessages() {
    for {
        // Grab the next message from the broadcast channel
        msg := <-broadcast

        // Send it out to every client that is currently connected
        for client := range clients {
            err := client.WriteJSON(msg)
            if err != nil {
                log.Printf("error: %v", err)
                client.Close()
                delete(clients, client)
            }
        }
    }
}

func semaphoreHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    defer conn.Close()

    // Register new client
    clients[conn] = true

    // Keep connection alive
    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            log.Printf("error: %v", err)
            delete(clients, conn)
            break
        }
    }
}

func main() {
    http.HandleFunc("/semaphore", semaphoreHandler)

    go handleMessages()

    // Simulate semaphore logic
    go func() {
        for {
            // Send green light
            broadcast <- map[string]string{"color": "green"}
            time.Sleep(5 * time.Second)
            // Send red light
            broadcast <- map[string]string{"color": "red"}
            time.Sleep(5 * time.Second)
        }
    }()

    log.Println("Server started on :8080")
    err := http.ListenAndServe(":8080", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
