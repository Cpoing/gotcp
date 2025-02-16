package chat

import (
	"fmt"
	"net"
	"sync"
)

type Manager struct {
	users map[string]net.Conn
	mutex sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		users: make(map[string]net.Conn),
	}
}

func (m *Manager) AddUser(username string, conn net.Conn) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.users[username] = conn

	m.broadcast("System", fmt.Sprintf("%s has joined the chat.\n", username))
}

func (m *Manager) RemoveUser(username string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.users, username)

	m.broadcast("System", fmt.Sprintf("%s has left the chat.\n", username))
}

func (m *Manager) Broadcast(sender, message string) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	formattedMessage := fmt.Sprintf("[%s]: %s", sender, message)
	for username, conn := range m.users {
		if username == sender {
			continue
		}

		_, err := fmt.Fprint(conn, formattedMessage)
		if err != nil {
			fmt.Printf("Error broadcasting to %s: %v\n", username, err)
		}
	}
}

func (m *Manager) broadcast(sender, message string) {
	formattedMessage := fmt.Sprintf("[%s]: %s", sender, message)
	for _, conn := range m.users {
		_, err := fmt.Fprint(conn, formattedMessage)
		if err != nil {
			fmt.Printf("Error broadcasting system message: %v\n", err)
		}
	}
}
