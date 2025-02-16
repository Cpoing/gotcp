package connection

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/Cpoing/gotcp/pkg/chat"
)

func HandleConnection(conn net.Conn, chatManager *chat.Manager) {
	defer conn.Close()

	conn.Write([]byte("Weolcome to ChatServer!\nPlease enter your username:\n"))

	reader := bufio.NewReader(conn)

	username, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Error reading username: %v", err)
		return
	}

	username = strings.TrimSpace(username)
	conn.Write([]byte(fmt.Sprintf("Hello %s, you have joined the chat!\n", username)))

	chatManager.AddUser(username, conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading message from %s: %v", username, err)
			break
		}

		trimmedMsg := strings.TrimSpace(message)
		if trimmedMsg == "/quit" {
			conn.Write([]byte("Goodbye!\n"))
			break
		}

		chatManager.Broadcast(username, message)
	}

	chatManager.RemoveUser(username)

}
