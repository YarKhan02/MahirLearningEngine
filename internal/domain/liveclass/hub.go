package liveclass

import (
	"context"
	"encoding/json"
	"sort"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

const writeTimeout = 10 * time.Second

// maxReadBytes lifts coder/websocket's default 32 KiB per-message read limit so
// image "files" messages (base64 data URLs) aren't rejected. Capped to bound abuse.
const maxReadBytes = 16 << 20 // 16 MB

// hostGraceDuration is how long a class stays live after the host's last
// connection drops, before it auto-ends. Absorbs brief reconnects (refresh,
// network blips, opening the board in a new tab).
const hostGraceDuration = 2 * time.Minute

// wsIn / wsOut are the whiteboard message envelopes.
type wsIn struct {
	Type  string          `json:"type"` // "scene-update" | "clear" | "files"
	Scene json.RawMessage `json:"scene,omitempty"`
	Files json.RawMessage `json:"files,omitempty"` // Excalidraw binary files map (images)
}

type wsOut struct {
	Type  string          `json:"type"` // "init" | "scene-update" | "clear" | "participants" | "files" | "ended"
	Scene json.RawMessage `json:"scene,omitempty"`
	Files json.RawMessage `json:"files,omitempty"`
	Count int             `json:"count,omitempty"`
	Names []string        `json:"names,omitempty"` // distinct student names present (for participants)
}

func mustJSON(v wsOut) []byte {
	b, _ := json.Marshal(v)
	return b
}

type client struct {
	conn   *websocket.Conn
	userID uuid.UUID
	name   string
	isHost bool
	send   chan []byte
	room   *room
}

type room struct {
	id      uuid.UUID
	clients map[*client]bool
	scene   json.RawMessage            // latest full scene (elements); replayed to joiners
	files   map[string]json.RawMessage // accumulated image files by id; replayed to joiners

	// Host presence tracking for auto-end (all guarded by mu).
	hostCount int
	hostID    uuid.UUID
	endTimer  *time.Timer

	mu sync.Mutex
}

func (r *room) setScene(scene json.RawMessage) {
	r.mu.Lock()
	r.scene = scene
	r.mu.Unlock()
}

// mergeFiles adds newly-uploaded image files to the room's store.
func (r *room) mergeFiles(raw json.RawMessage) {
	if len(raw) == 0 {
		return
	}
	var incoming map[string]json.RawMessage
	if json.Unmarshal(raw, &incoming) != nil {
		return
	}
	r.mu.Lock()
	for k, v := range incoming {
		r.files[k] = v
	}
	r.mu.Unlock()
}

// filesJSON marshals the accumulated files, or nil when empty.
func (r *room) filesJSON() json.RawMessage {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.files) == 0 {
		return nil
	}
	b, _ := json.Marshal(r.files)
	return b
}

func (r *room) clearFiles() {
	r.mu.Lock()
	r.files = make(map[string]json.RawMessage)
	r.mu.Unlock()
}

func (r *room) add(c *client) {
	r.mu.Lock()
	r.clients[c] = true
	r.mu.Unlock()
}

func (r *room) remove(c *client) {
	r.mu.Lock()
	delete(r.clients, c)
	r.mu.Unlock()
}

// broadcast pushes msg to every client (optionally excluding one). A full send
// buffer means a slow client; we drop the frame — scene-updates are full-state
// so the next one supersedes it, and joiners always get the latest via "init".
func (r *room) broadcast(msg []byte, except *client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for c := range r.clients {
		if c == except {
			continue
		}
		select {
		case c.send <- msg:
		default:
		}
	}
}

// broadcastPresence sends the distinct participant count plus the de-duplicated
// list of student (non-host) names currently connected.
func (r *room) broadcastPresence() {
	r.mu.Lock()
	users := make(map[uuid.UUID]bool)
	seenNames := make(map[uuid.UUID]bool)
	names := make([]string, 0)
	for c := range r.clients {
		users[c.userID] = true
		if !c.isHost && !seenNames[c.userID] {
			seenNames[c.userID] = true
			if c.name != "" {
				names = append(names, c.name)
			}
		}
	}
	r.mu.Unlock()
	sort.Strings(names)
	r.broadcast(mustJSON(wsOut{Type: "participants", Count: len(users), Names: names}), nil)
}

// Hub keeps one in-memory room per live session. The single-instance design is
// deliberate; a Redis Pub/Sub fan-out can replace broadcast() when scaling out.
type Hub struct {
	rooms map[uuid.UUID]*room
	mu    sync.Mutex
	// endFn ends a session in the DB when the host is gone past the grace period.
	endFn func(sessionID, hostID uuid.UUID)
	grace time.Duration
}

func NewHub(endFn func(sessionID, hostID uuid.UUID)) *Hub {
	return &Hub{
		rooms: make(map[uuid.UUID]*room),
		endFn: endFn,
		grace: hostGraceDuration,
	}
}

func (h *Hub) getRoom(id uuid.UUID) *room {
	h.mu.Lock()
	defer h.mu.Unlock()
	r, ok := h.rooms[id]
	if !ok {
		r = &room{
			id:      id,
			clients: make(map[*client]bool),
			files:   make(map[string]json.RawMessage),
		}
		h.rooms[id] = r
	}
	return r
}

func (h *Hub) dropIfEmpty(r *room) {
	h.mu.Lock()
	defer h.mu.Unlock()
	r.mu.Lock()
	// Keep the room while an auto-end timer is pending, so a host reconnect
	// reuses it and can cancel the timer.
	empty := len(r.clients) == 0 && r.endTimer == nil
	r.mu.Unlock()
	if empty {
		delete(h.rooms, r.id)
	}
}

// hostJoined records a host connection and cancels any pending auto-end.
func (h *Hub) hostJoined(r *room, hostID uuid.UUID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.hostCount++
	r.hostID = hostID
	if r.endTimer != nil {
		r.endTimer.Stop()
		r.endTimer = nil
	}
}

// hostLeft decrements host presence and, when the last host is gone, arms the
// grace timer that auto-ends the class.
func (h *Hub) hostLeft(r *room) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.hostCount--
	if r.hostCount > 0 {
		return
	}
	r.hostCount = 0
	sid, hid := r.id, r.hostID
	r.endTimer = time.AfterFunc(h.grace, func() {
		r.mu.Lock()
		stillGone := r.hostCount == 0
		r.endTimer = nil
		r.mu.Unlock()
		if !stillGone {
			return
		}
		if h.endFn != nil {
			h.endFn(sid, hid)
		}
		r.broadcast(mustJSON(wsOut{Type: "ended"}), nil)
		h.dropIfEmpty(r)
	})
}

// Serve runs one connection to completion (blocks until disconnect).
func (h *Hub) Serve(ctx context.Context, conn *websocket.Conn, sessionID, userID uuid.UUID, name string, isHost bool) {
	conn.SetReadLimit(maxReadBytes) // allow large image (files) messages
	r := h.getRoom(sessionID)
	c := &client{conn: conn, userID: userID, name: name, isHost: isHost, send: make(chan []byte, 32), room: r}
	r.add(c)
	if isHost {
		h.hostJoined(r, userID)
	}

	// Replay current files then the scene, so images resolve when the elements
	// that reference them arrive (Excalidraw needs the files registered first).
	r.mu.Lock()
	scene := r.scene
	r.mu.Unlock()
	if f := r.filesJSON(); f != nil {
		c.send <- mustJSON(wsOut{Type: "files", Files: f})
	}
	c.send <- mustJSON(wsOut{Type: "init", Scene: scene})
	r.broadcastPresence()

	writerDone := make(chan struct{})
	go c.writeLoop(ctx, writerDone)

	c.readLoop(ctx) // blocks until the client goes away

	r.remove(c)
	if isHost {
		h.hostLeft(r) // arms the grace timer when the last host connection drops
	}
	close(c.send)
	<-writerDone
	r.broadcastPresence()
	h.dropIfEmpty(r)
	_ = conn.Close(websocket.StatusNormalClosure, "")
}

func (c *client) writeLoop(ctx context.Context, done chan struct{}) {
	defer close(done)
	for msg := range c.send {
		wctx, cancel := context.WithTimeout(ctx, writeTimeout)
		err := c.conn.Write(wctx, websocket.MessageText, msg)
		cancel()
		if err != nil {
			return
		}
	}
}

func (c *client) readLoop(ctx context.Context) {
	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			return
		}
		var in wsIn
		if json.Unmarshal(data, &in) != nil {
			continue
		}
		switch in.Type {
		case "scene-update":
			if !c.isHost {
				continue // students cannot draw, whatever their client sends
			}
			c.room.setScene(in.Scene)
			c.room.broadcast(mustJSON(wsOut{Type: "scene-update", Scene: in.Scene}), c)
		case "files":
			if !c.isHost {
				continue // students cannot add images
			}
			c.room.mergeFiles(in.Files)
			c.room.broadcast(mustJSON(wsOut{Type: "files", Files: in.Files}), c)
		case "clear":
			if !c.isHost {
				continue
			}
			c.room.setScene(nil)
			c.room.clearFiles()
			c.room.broadcast(mustJSON(wsOut{Type: "clear"}), nil)
		}
	}
}
