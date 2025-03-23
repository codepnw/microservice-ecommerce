package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/codepnw/microservice-ecommerce/ecom-api/store"
	"github.com/codepnw/microservice-ecommerce/token"
	"github.com/codepnw/microservice-ecommerce/utils"
	"github.com/gin-gonic/gin"
)

func (h *handler) createUser(c *gin.Context) {
	var u UserReq

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// hash password
	hashed, err := utils.HashPassword(u.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	u.Password = hashed

	created, err := h.server.CreateUser(c.Request.Context(), toStoreUser(u))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := toUserRes(created)
	c.JSON(http.StatusCreated, res)
}

func (h *handler) listUsers(c *gin.Context) {
	users, err := h.server.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var res ListUserRes
	for _, u := range users {
		res.Users = append(res.Users, toUserRes(&u))
	}

	c.JSON(http.StatusOK, res)
}

func (h *handler) updateUser(c *gin.Context) {
	var u UserReq
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get Context
	claims, exists := c.Get(claimsKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	email := claims.(*token.UserClaims).Email

	user, err := h.server.GetUser(c.Request.Context(), email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// patch user req
	patchUserReq(user, u)
	if user.Email == "" {
		user.Email = email
	}

	updated, err := h.server.UpdateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := toUserRes(updated)
	c.JSON(http.StatusOK, res)
}

func (h *handler) deleteUser(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error pasing ID"})
		return
	}

	if err := h.server.DeleteUser(c.Request.Context(), idInt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *handler) loginUser(c *gin.Context) {
	var u LoginUserReq

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	gu, err := h.server.GetUser(c.Request.Context(), u.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := utils.CheckPassword(u.Password, gu.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong email or password"})
		return
	}

	// create JWT
	accessToken, accessClaims, err := h.TokenMaker.CreateToken(gu.ID, gu.Email, gu.IsAdmin, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	refreshToken, refreshClaims, err := h.TokenMaker.CreateToken(gu.ID, gu.Email, gu.IsAdmin, 24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println(refreshClaims.RegisteredClaims.ExpiresAt.Time)

	// Session
	session, err := h.server.CreateSession(c.Request.Context(), &store.Session{
		ID:           refreshClaims.RegisteredClaims.ID,
		UserEmail:    gu.Email,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error-3": err.Error()})
		return
	}

	res := LoginUserRes{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.RegisteredClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.RegisteredClaims.ExpiresAt.Time,
		User:                  toUserRes(gu),
	}

	c.JSON(http.StatusOK, res)
}

func (h *handler) logoutUser(c *gin.Context) {
	// Get Context
	claims, exists := c.Get(claimsKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	id := claims.(*token.UserClaims).RegisteredClaims.ID

	if err := h.server.DeleteSession(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *handler) renewAccessToken(c *gin.Context) {
	var req RenewAccessTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshClaims, err := h.TokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	session, err := h.server.GetSession(c.Request.Context(), refreshClaims.RegisteredClaims.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if session.IsRevoked {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "session revoked"})
		return
	}

	if session.UserEmail != refreshClaims.Email {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	accessToken, accessClaims, err := h.TokenMaker.CreateToken(refreshClaims.ID, refreshClaims.Email, refreshClaims.IsAdmin, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := RenewAccessTokenRes{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.RegisteredClaims.ExpiresAt.Time,
	}

	c.JSON(http.StatusOK, res)
}

func (h *handler) revokeSession(c *gin.Context) {
	// Get Context
	claims, exists := c.Get(claimsKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	id := claims.(*token.UserClaims).RegisteredClaims.ID

	if err := h.server.RevokeSession(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func toStoreUser(u UserReq) *store.User {
	return &store.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: u.Password,
		IsAdmin:  u.IsAdmin,
	}
}

func toUserRes(u *store.User) UserRes {
	return UserRes{
		Name:    u.Name,
		Email:   u.Email,
		IsAdmin: u.IsAdmin,
	}
}

func patchUserReq(user *store.User, u UserReq) {
	if u.Name != "" {
		user.Name = u.Name
	}
	if u.Email != "" {
		user.Email = u.Email
	}
	if u.Password != "" {
		hashed, err := utils.HashPassword(u.Password)
		if err != nil {
			panic(err)
		}
		user.Password = hashed
	}
	if u.IsAdmin {
		user.IsAdmin = u.IsAdmin
	}
	user.UpdatedAt = toTimePtr(time.Now())
}
