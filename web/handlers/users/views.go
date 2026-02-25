package users

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Doris-Mwito5/ginja-ai/internal/apperr"
	"github.com/Doris-Mwito5/ginja-ai/internal/db"
	"github.com/Doris-Mwito5/ginja-ai/internal/dtos"
	"github.com/Doris-Mwito5/ginja-ai/internal/jwt"
	"github.com/Doris-Mwito5/ginja-ai/internal/services"
	"github.com/Doris-Mwito5/ginja-ai/internal/utils"
	"github.com/gin-gonic/gin"
)

const handlerTimeout = 15 * time.Second

func register(
	dB db.DB,
	userService services.UserService,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req dtos.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.HandleError(c, apperr.NewErrorWithType(err, apperr.BadRequest))
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), handlerTimeout)
		defer cancel()

		user, err := userService.Register(ctx, dB, &req)
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				utils.HandleError(c, apperr.NewInternal("request timed out; check database is running and reachable"))
				return
			}
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

func login(
	dB db.DB,
	userService services.UserService,
	jwtMaker jwt.JWTToken,
	tokenDuration time.Duration,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req dtos.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.HandleError(c, apperr.NewErrorWithType(err, apperr.BadRequest))
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), handlerTimeout)
		defer cancel()

		user, err := userService.ValidateCredentials(ctx, dB, strings.TrimSpace(req.Username), req.Password)
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				utils.HandleError(c, apperr.NewInternal("request timed out; check database is running and reachable"))
				return
			}
			utils.HandleError(c, err)
			return
		}

		accessToken, err := jwtMaker.CreateToken(user.Username, tokenDuration)
		if err != nil {
			utils.HandleError(c, apperr.NewInternal("failed to create token"))
			return
		}

		expiresAt := time.Now().Add(tokenDuration).Format(time.RFC3339)
		c.JSON(http.StatusOK, dtos.LoginResponse{
			AccessToken: accessToken,
			User:        user,
			ExpiresAt:   expiresAt,
		})
	}
}

func getUserByID(
	dB db.DB,
	userService services.UserService,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			utils.HandleError(c, apperr.NewBadRequest("invalid user id"))
			return
		}

		user, err := userService.GetUserByID(c.Request.Context(), dB, id)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func getUserByUsername(
	dB db.DB,
	userService services.UserService,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		username := strings.TrimSpace(c.Param("username"))
		if username == "" {
			utils.HandleError(c, apperr.NewBadRequest("username is required"))
			return
		}

		user, err := userService.GetUserByUsername(c.Request.Context(), dB, username)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
