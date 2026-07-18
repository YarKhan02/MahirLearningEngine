package user

import (
	"net"
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/middleware"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/response"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/common"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/logging"
	"github.com/YarKhan02/MahirLearningEngine/internal/infrastructure/metrics"
	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	svc           *Service
	studentSvc    common.StudentProfileProvider
	tokenSvc      *token.Service
	secureCookies bool
}

func NewHandler(svc *Service, studentSvc common.StudentProfileProvider, tokenSvc *token.Service, secureCookies bool) *Handler {
	return &Handler{
		svc:           svc,
		studentSvc:    studentSvc,
		tokenSvc:      tokenSvc,
		secureCookies: secureCookies,
	}
}

const refreshCookieName = "refresh_token"
const refreshCookieMaxAge = 30 * 24 * 60 * 60 // 30 days, in seconds

// setRefreshCookie writes the HttpOnly refresh-token cookie with attributes
// appropriate to the environment. Over HTTPS (production) it uses
// SameSite=None + Secure so the cookie survives a cross-site request from the
// SPA; over plain HTTP (local dev) those attributes would make the browser drop
// the cookie, so it falls back to SameSite=Lax without Secure.
func (h *Handler) setRefreshCookie(c *gin.Context, value string, maxAge int) {
	cookie := &http.Cookie{
		Name:     refreshCookieName,
		Value:    value,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   maxAge,
		Secure:   h.secureCookies,
		SameSite: http.SameSiteLaxMode,
	}
	if h.secureCookies {
		cookie.SameSite = http.SameSiteNoneMode
	}
	http.SetCookie(c.Writer, cookie)
}

func (h *Handler) RegisterAdmin(c *gin.Context) {

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := req.Validate(); err != nil {
		response.WriteError(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	u, err := h.svc.RegisterAdmin(c.Request.Context(), req)
	if err != nil {
		switch err {
		case ErrEmailTaken:
			response.WriteError(c, http.StatusConflict, "email already registered")
		default:
			response.WriteInternal(c, err)
		}
		return
	}

	response.WriteJSON(c, http.StatusCreated, UserResponse{ID: u.ID, Email: u.Email, Role: u.Role})
}

func (h *Handler) Login(c *gin.Context) {

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := req.Validate(); err != nil {
		response.WriteError(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	u, err := h.svc.Authenticate(c.Request.Context(), req.Identifier, req.Password)
	if err != nil {
		reason := "invalid_credentials"
		switch err {
		case ErrAccountLocked:
			reason = "account_locked"
			response.WriteError(c, http.StatusLocked, "account temporarily locked")
		case ErrAccountBanned:
			reason = "account_banned"
			response.WriteError(c, http.StatusForbidden, "account banned")
		default:
			// always same message — don't leak which field is wrong
			response.WriteError(c, http.StatusUnauthorized, "invalid credentials")
		}
		metrics.RecordSecurityEvent("login_failed", reason)
		metrics.RecordLogin("failure")
		logging.FromLogger(c.Request.Context()).Warn("login failed",
			zap.String("event", "login_failed"),
			zap.String("identifier", req.Identifier),
			zap.String("reason", reason),
		)
		return
	}

	metrics.RecordLogin("success")

	accessToken, err := h.tokenSvc.IssueAccessToken(u.ID, u.Email, u.Role)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	ip := c.Request.RemoteAddr
	if host, _, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
		ip = host
	}

	refreshToken, err := h.tokenSvc.IssueRefreshToken(c.Request.Context(), u.ID, ip, c.Request.UserAgent())
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	// refresh token in HttpOnly cookie, access token in body
	h.setRefreshCookie(c, refreshToken, refreshCookieMaxAge)

	// Admins carry their email on the account; students don't, so pull their
	// real name and contact email from the student record for portal display.
	authUser := AuthUser{
		ID:       u.ID.String(),
		Name:     u.Email,
		Username: u.Username,
		Email:    u.Email,
		Role:     u.Role,
	}
	if u.Role == "student" {
		if profile, perr := h.studentSvc.GetStudentProfile(c.Request.Context(), u.ID); perr == nil && profile != nil {
			authUser.Name = profile.FullName
			authUser.Email = profile.Email
			authUser.Username = profile.Username
		}
	}

	response.WriteJSON(c, http.StatusOK, LoginResponse{
		AccessToken: accessToken,
		User:        authUser,
	})
}

func (h *Handler) Refresh(c *gin.Context) {

	cookie, err := c.Request.Cookie("refresh_token")
	if err != nil {
		response.WriteError(c, http.StatusUnauthorized, "missing refresh token")
		return
	}

	rt, newRawToken, err := h.tokenSvc.RotateRefreshToken(c.Request.Context(), cookie.Value)
	if err != nil {
		response.WriteError(c, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	u, err := h.svc.FindByID(c.Request.Context(), rt.UserID)
	if err != nil {
		response.WriteError(c, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	accessToken, err := h.tokenSvc.IssueAccessToken(u.ID, u.Email, u.Role)
	if err != nil {
		response.WriteInternal(c, err)
		return
	}

	h.setRefreshCookie(c, newRawToken, refreshCookieMaxAge)

	response.WriteJSON(c, http.StatusOK, TokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   900,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	cookie, err := c.Request.Cookie(refreshCookieName)
	if err == nil {
		_ = h.tokenSvc.RevokeByRawToken(c.Request.Context(), cookie.Value)
	}

	h.setRefreshCookie(c, "", -1)

	response.WriteJSON(c, http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}

func (h *Handler) ResetPassword(c *gin.Context) {

	userID, err := uuid.Parse(c.Param("userID"))
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid user id")
		return
	}

	currentUserID, check := middleware.CurrentUserID(c)
	if !check {
		response.WriteError(c, http.StatusBadRequest, "invaid request")
		return
	}
	
	if currentUserID != userID {
		response.WriteError(c, http.StatusBadRequest, "invaid request")
		return
	}


	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	err = h.svc.ResetPassword(c.Request.Context(), userID, req.NewPassword)
	if err != nil {
		response.WriteError(c, http.StatusBadRequest, "reset password failed")
		return
	}

	response.WriteJSON(c, http.StatusOK, gin.H{
		"message": "password reset successful",
	})
}