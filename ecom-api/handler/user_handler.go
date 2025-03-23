package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/codepnw/microservice-ecommerce/ecom-api/store"
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

	user, err := h.server.GetUser(c.Request.Context(), u.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// patch user req
	patchUserReq(user, u)

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
	accessToken, _, err := h.tokenMaker.CreateToken(gu.ID, gu.Email, gu.IsAdmin, 15*time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := LoginUserRes{
		AccessToken: accessToken,
		User: UserRes{
			Name:    gu.Name,
			Email:   gu.Email,
			IsAdmin: gu.IsAdmin,
		},
	}

	c.JSON(http.StatusOK, res)
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
