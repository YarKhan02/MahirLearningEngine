package handler

import (
	"net"
	"net/http"

	"github.com/YarKhan02/MahirLearningEngine/internal/api/dto"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/token"
	"github.com/YarKhan02/MahirLearningEngine/internal/domain/user"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userSvc  *user.Service
	tokenSvc *token.Service
}

func NewAuthHandler(userSvc *user.Service, tokenSvc *token.Service) *AuthHandler {
	return &AuthHandler{userSvc: userSvc, tokenSvc: tokenSvc}
}

func (h *AuthHandler) RegisterAdmin(c *gin.Context) {
	
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	if err := req.Validate(); err != nil {
		writeError(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	u, err := h.userSvc.RegisterAdmin(c.Request.Context(), req)
	if err != nil {
		switch err {
		case user.ErrEmailTaken:
			writeError(c, http.StatusConflict, "email already registered")
		default:
			writeError(c, http.StatusInternalServerError, "registration failed")
		}
		return
	}

	writeJSON(c, http.StatusCreated, dto.UserResponse{ID: u.ID, Email: u.Email, Role: u.Role})
}

func (h *AuthHandler) Login(c *gin.Context) {

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	if err := req.Validate(); err != nil {
		writeError(c, http.StatusUnprocessableEntity, err.Error())
		return
	}

	u, err := h.userSvc.Authenticate(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case user.ErrAccountLocked:
			writeError(c, http.StatusLocked, "account temporarily locked")
		case user.ErrAccountBanned:
			writeError(c, http.StatusForbidden, "account banned")
		default:
			// always same message — don't leak which field is wrong
			writeError(c, http.StatusUnauthorized, "invalid credentials")
		}
		return
	}

	accessToken, err := h.tokenSvc.IssueAccessToken(u)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to issue access token")
		return
	}

	ip := c.Request.RemoteAddr
	if host, _, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
		ip = host
	}
	
	refreshToken, err := h.tokenSvc.IssueRefreshToken(c, u, ip, c.Request.UserAgent())
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to issue refresh token")
		return
	}

	// refresh token in HttpOnly cookie, access token in body
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60,
	})

	writeJSON(c, http.StatusOK, dto.LoginResponse{
		AccessToken: accessToken,
		User: dto.AuthUser{
			ID: 	u.ID.String(),
			Name: 	"yarkhan",
			Email: 	u.Email,
			Role: 	u.Role,
		},
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	
	cookie, err := c.Request.Cookie("refresh_token")
	if err != nil {
		writeError(c, http.StatusUnauthorized, "missing refresh token")
		return
	}

	rt, newRawToken, err := h.tokenSvc.RotateRefreshToken(c.Request.Context(), cookie.Value)
	if err != nil {
		writeError(c, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	u, err := h.userSvc.FindByID(c.Request.Context(), rt.UserID)
	if err != nil {
		writeError(c, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	accessToken, err := h.tokenSvc.IssueAccessToken(u)
	if err != nil {
		writeError(c, http.StatusInternalServerError, "failed to issue access token")
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRawToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60,
	})

	writeJSON(c, http.StatusOK, dto.TokenResponse{
		AccessToken: accessToken,
		ExpiresIn:   900,
	})
}