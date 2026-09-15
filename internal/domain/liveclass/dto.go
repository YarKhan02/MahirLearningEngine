package liveclass

import "encoding/json"

type StartSessionRequest struct {
	BatchID  string `json:"batchId"`
	CourseID string `json:"courseId"`
	Title    string `json:"title"`
}

type SaveSnapshotRequest struct {
	Scene    json.RawMessage `json:"scene"`
	ImageURL string          `json:"imageUrl"`
}

type SessionResponse struct {
	ID          string `json:"id"`
	BatchID     string `json:"batchId"`
	CourseID    string `json:"courseId"`
	HostID      string `json:"hostId"`
	Title       string `json:"title"`
	ClassDate   string `json:"classDate"` // YYYY-MM-DD
	Status          string `json:"status"`
	StartedAt       string `json:"startedAt"`
	EndedAt         string `json:"endedAt,omitempty"`
	DurationSeconds int64  `json:"durationSeconds"` // 0 while live / not ended
	CourseTitle     string `json:"courseTitle,omitempty"`
	BatchName       string `json:"batchName,omitempty"`
}

type BatchOptionResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CourseOptionResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// TicketResponse is returned to the client to open the WebSocket.
type TicketResponse struct {
	Ticket string `json:"ticket"`
}

// RTCTokenResponse carries the LiveKit server URL and a join token.
type RTCTokenResponse struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

// MicRequest lets the host grant/revoke a student's ability to speak.
type MicRequest struct {
	Identity string `json:"identity"` // the student's LiveKit identity (their user id)
	Allow    bool   `json:"allow"`
}

type SnapshotResponse struct {
	ID        string `json:"id"`
	SessionID string `json:"sessionId"`
	ImageURL  string `json:"imageUrl"`
	Version   int    `json:"version"`
	CreatedAt string `json:"createdAt"`
}
