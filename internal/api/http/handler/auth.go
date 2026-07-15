package handler

import (
	"net"
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/api/http/response"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/student"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/user"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userSvc    *user.Service
	studentSvc *student.Service
	tokenSvc   *token.Service
	secureCookies bool
}

func NewAuthHandler(userSvc *user.Service, studentSvc *student.Service, tokenSvc *token.Service, secureCookies bool) *AuthHandler {
	return &AuthHandler{userSvc: userSvc, studentSvc: studentSvc, tokenSvc: tokenSvc, secureCookies: secureCookies}
}

const refreshCookieName = "refresh_token"
const refreshCookieMaxAge = 30 * 24 * 60 * 60 // 30 days, in seconds

// setRefreshCookie writes the HttpOnly refresh-token cookie with attributes
// appropriate to the environment. Over HTTPS (production) it uses
// SameSite=None + Secure so the cookie survives a cross-site request from the
// SPA; over plain HTTP (local dev) those attributes would make the browser drop
// the cookie, so it falls back to SameSite=Lax without Secure.
func (h *AuthHandler) setRefreshCookie(c *gin.Context, value string, maxAge int) {
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

func (h *AuthHandler) RegisterAdmin(c *gin.Context) {
	
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := req.Validate(); err != nil {
		response.WriteError(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	u, err := h.userSvc.RegisterAdmin(c.Request.Context(), req)
	if err != nil {
		switch err {
		case user.ErrEmailTaken:
			response.WriteError(c, http.StatusConflict, "email already registered")
		default:
			response.WriteError(c, http.StatusInternalServerError, "registration failed")
		}
		return
	}

	response.WriteJSON(c, http.StatusCreated, dto.UserResponse{ID: u.ID, Email: u.Email, Role: u.Role})
}

func (h *AuthHandler) Login(c *gin.Context) {

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.WriteError(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := req.Validate(); err != nil {
		response.WriteError(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	u, err := h.userSvc.Authenticate(c.Request.Context(), req.Identifier, req.Password)
	if err != nil {
		switch err {
		case user.ErrAccountLocked:
			response.WriteError(c, http.StatusLocked, "account temporarily locked")
		case user.ErrAccountBanned:
			response.WriteError(c, http.StatusForbidden, "account banned")
		default:
			// always same message — don't leak which field is wrong
			response.WriteError(c, http.StatusUnauthorized, "invalid credentials")
		}
		return
	}

	accessToken, err := h.tokenSvc.IssueAccessToken(u)
	if err != nil {
		response.WriteError(c, http.StatusInternalServerError, "failed to issue access token")
		return
	}

	ip := c.Request.RemoteAddr
	if host, _, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
		ip = host
	}
	
	refreshToken, err := h.tokenSvc.IssueRefreshToken(c, u, ip, c.Request.UserAgent())
	if err != nil {
		response.WriteError(c, http.StatusInternalServerError, "failed to issue refresh token")
		return
	}

	// refresh token in HttpOnly cookie, access token in body
	h.setRefreshCookie(c, refreshToken, refreshCookieMaxAge)

	// Admins carry their email on the account; students don't, so pull their
	// real name and contact email from the student record for portal display.
	authUser := dto.AuthUser{
		ID:       u.ID.String(),
		Name:     u.Email,
		Username: u.Username,
		Email:    u.Email,
		Role:     u.Role,
	}
	if u.Role == "student" {
		if profile, perr := h.studentSvc.GetProfileByUserID(c.Request.Context(), u.ID); perr == nil && profile != nil {
			authUser.Name = profile.FullName
			authUser.Email = profile.Email
			authUser.Username = profile.Username
		}
	}

	response.WriteJSON(c, http.StatusOK, dto.LoginResponse{
		AccessToken: accessToken,
		User:        authUser,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	
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

	u, err := h.userSvc.FindByID(c.Request.Context(), rt.UserID)
	if err != nil {
		response.WriteError(c, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	accessToken, err := h.tokenSvc.IssueAccessToken(u)
	if err != nil {
		response.WriteError(c, http.StatusInternalServerError, "failed to issue access token")
		return
	}

	h.setRefreshCookie(c, newRawToken, refreshCookieMaxAge)

	response.WriteJSON(c, http.StatusOK, dto.TokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   900,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	cookie, err := c.Request.Cookie(refreshCookieName)
	if err == nil {
		_ = h.tokenSvc.RevokeByRawToken(c.Request.Context(), cookie.Value)
	}

	h.setRefreshCookie(c, "", -1)

	response.WriteJSON(c, http.StatusOK, gin.H{
		"message": "logged out successfully",
	})
}