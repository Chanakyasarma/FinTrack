package websocket_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"fintrack/internal/websocket"

	gorillaws "github.com/gorilla/websocket"
)

var upgrader = gorillaws.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// startTestServer creates an HTTP test server that upgrades to WS
// and registers the client with the hub.
func startTestServer(t *testing.T, hub *websocket.Hub, userID string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("upgrade error: %v", err)
			return
		}
		client := websocket.NewClient(hub, conn, userID)
		hub.Register(client)
		go client.WritePump()
		go client.ReadPump()
	}))
	return srv
}

func dial(t *testing.T, srv *httptest.Server) *gorillaws.Conn {
	t.Helper()
	url := "ws" + strings.TrimPrefix(srv.URL, "http")
	conn, _, err := gorillaws.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	return conn
}

func TestHub_BroadcastToUser(t *testing.T) {
	hub := websocket.NewHub()
	go hub.Run()
	time.Sleep(10 * time.Millisecond) // let hub goroutine start

	srv := startTestServer(t, hub, "user-1")
	defer srv.Close()

	conn := dial(t, srv)
	defer conn.Close()
	time.Sleep(30 * time.Millisecond) // let client register

	// Send an event to user-1
	hub.SendToUser("user-1", websocket.Event{
		Type:    websocket.EventTransactionCreated,
		Payload: map[string]any{"amount": 42.0},
	})

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("reading message: %v", err)
	}

	var event websocket.Event
	if err := json.Unmarshal(msg, &event); err != nil {
		t.Fatalf("unmarshaling event: %v", err)
	}
	if event.Type != websocket.EventTransactionCreated {
		t.Errorf("expected type %q, got %q", websocket.EventTransactionCreated, event.Type)
	}
}

func TestHub_DoesNotBroadcastToOtherUser(t *testing.T) {
	hub := websocket.NewHub()
	go hub.Run()
	time.Sleep(10 * time.Millisecond)

	srv1 := startTestServer(t, hub, "user-1")
	defer srv1.Close()
	srv2 := startTestServer(t, hub, "user-2")
	defer srv2.Close()

	conn1 := dial(t, srv1)
	defer conn1.Close()
	conn2 := dial(t, srv2)
	defer conn2.Close()
	time.Sleep(30 * time.Millisecond)

	// Send only to user-2
	hub.SendToUser("user-2", websocket.Event{
		Type:    websocket.EventBalanceUpdated,
		Payload: map[string]any{"balance": 9999.0},
	})

	// conn2 should receive the message
	conn2.SetReadDeadline(time.Now().Add(time.Second))
	_, msg, err := conn2.ReadMessage()
	if err != nil {
		t.Fatalf("user-2 should have received the message: %v", err)
	}
	var event websocket.Event
	json.Unmarshal(msg, &event)
	if event.Type != websocket.EventBalanceUpdated {
		t.Errorf("user-2: unexpected event type %q", event.Type)
	}

	// conn1 should NOT receive anything within the timeout
	conn1.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	_, _, err = conn1.ReadMessage()
	if err == nil {
		t.Error("user-1 should NOT have received user-2's message")
	}
}

func TestHub_MultipleClientsForSameUser(t *testing.T) {
	hub := websocket.NewHub()
	go hub.Run()
	time.Sleep(10 * time.Millisecond)

	// Two browser tabs, same user
	srv := startTestServer(t, hub, "user-multi")
	defer srv.Close()

	conn1 := dial(t, srv)
	defer conn1.Close()
	conn2 := dial(t, srv)
	defer conn2.Close()
	time.Sleep(50 * time.Millisecond)

	hub.SendToUser("user-multi", websocket.Event{
		Type:    websocket.EventTransactionCreated,
		Payload: map[string]any{"id": "tx-1"},
	})

	var wg sync.WaitGroup
	received := make([]bool, 2)

	for i, c := range []*gorillaws.Conn{conn1, conn2} {
		wg.Add(1)
		go func(idx int, conn *gorillaws.Conn) {
			defer wg.Done()
			conn.SetReadDeadline(time.Now().Add(time.Second))
			_, _, err := conn.ReadMessage()
			if err == nil {
				received[idx] = true
			}
		}(i, c)
	}

	wg.Wait()

	if !received[0] {
		t.Error("first connection (tab 1) should have received the broadcast")
	}
	if !received[1] {
		t.Error("second connection (tab 2) should have received the broadcast")
	}
}
