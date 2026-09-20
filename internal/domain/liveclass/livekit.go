package liveclass

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var ErrVideoNotConfigured = errors.New("video is not configured")

// LiveKit mints join tokens and manages participant permissions. It talks to
// LiveKit directly (JWT + RoomService HTTP/JSON) to avoid the heavy official
// Go SDK dependency tree — a LiveKit token is just an HS256 JWT.
type LiveKit struct {
	wsURL     string
	hostURL   string
	apiKey    string
	apiSecret string
	http      *http.Client
}

func NewLiveKit(wsURL, hostURL, apiKey, apiSecret string) *LiveKit {
	return &LiveKit{
		wsURL:     wsURL,
		hostURL:   hostURL,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		http:      &http.Client{Timeout: 10 * time.Second},
	}
}

func (lk *LiveKit) Configured() bool {
	return lk != nil && lk.apiKey != "" && lk.apiSecret != "" && lk.wsURL != ""
}

func (lk *LiveKit) WSURL() string { return lk.wsURL }

// videoGrant is LiveKit's "video" claim.
type videoGrant struct {
	RoomJoin       bool   `json:"roomJoin,omitempty"`
	RoomAdmin      bool   `json:"roomAdmin,omitempty"`
	RoomCreate     bool   `json:"roomCreate,omitempty"`
	Room           string `json:"room,omitempty"`
	CanPublish     *bool  `json:"canPublish,omitempty"`
	CanSubscribe   *bool  `json:"canSubscribe,omitempty"`
	CanPublishData *bool  `json:"canPublishData,omitempty"`
}

func (lk *LiveKit) signToken(identity, name string, grant videoGrant, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"iss":   lk.apiKey,
		"nbf":   now.Unix(),
		"exp":   now.Add(ttl).Unix(),
		"video": grant,
	}
	if identity != "" {
		claims["sub"] = identity
	}
	if name != "" {
		claims["name"] = name
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(lk.apiSecret))
}

// Token mints a room-join token. canPublish grants camera/mic/screen publishing
// (true for the host; false for students until they're allowed to speak).
func (lk *LiveKit) Token(identity, name, room string, canPublish bool) (string, error) {
	if !lk.Configured() {
		return "", ErrVideoNotConfigured
	}
	pub := canPublish
	sub, data := true, true
	return lk.signToken(identity, name, videoGrant{
		RoomJoin:       true,
		Room:           room,
		CanPublish:     &pub,
		CanSubscribe:   &sub,
		CanPublishData: &data,
	}, 2*time.Hour)
}

// SetPublish grants/revokes a connected participant's publish permission via the
// LiveKit RoomService (Twirp/JSON) — used for raise-hand. When granting, publish
// is scoped to MICROPHONE only: "allowed to speak" must not also let a student
// turn on their camera or share their screen.
func (lk *LiveKit) SetPublish(ctx context.Context, room, identity string, allow bool) error {
	if !lk.Configured() || lk.hostURL == "" {
		return ErrVideoNotConfigured
	}
	perm := map[string]any{
		"canSubscribe":   true,
		"canPublish":     allow,
		"canPublishData": true,
	}
	if allow {
		// Restrict the grant to audio so a raised hand can't publish video/screen.
		perm["canPublishSources"] = []string{"MICROPHONE"}
	}
	return lk.roomService(ctx, "UpdateParticipant",
		videoGrant{RoomAdmin: true, Room: room},
		map[string]any{"room": room, "identity": identity, "permission": perm})
}

// DeleteRoom disconnects everyone from a room and tears it down. Called when a
// class ends so lingering join tokens (which LiveKit does not re-check against
// our DB) can't keep a removed participant connected past the class.
func (lk *LiveKit) DeleteRoom(ctx context.Context, room string) error {
	if !lk.Configured() || lk.hostURL == "" {
		return ErrVideoNotConfigured
	}
	return lk.roomService(ctx, "DeleteRoom",
		videoGrant{RoomCreate: true, Room: room},
		map[string]any{"room": room})
}

// roomService POSTs a Twirp/JSON request to LiveKit's RoomService, signing an
// admin token with the given grant.
func (lk *LiveKit) roomService(ctx context.Context, method string, grant videoGrant, payload map[string]any) error {
	adminTok, err := lk.signToken("", "", grant, 10*time.Minute)
	if err != nil {
		return err
	}
	body, _ := json.Marshal(payload)

	url := strings.TrimRight(lk.hostURL, "/") + "/twirp/livekit.RoomService/" + method
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminTok)

	resp, err := lk.http.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("livekit %s: %s: %s", method, resp.Status, strings.TrimSpace(string(b)))
	}
	return nil
}
