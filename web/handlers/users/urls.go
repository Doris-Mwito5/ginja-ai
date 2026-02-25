package users

import (
	"github.com/Doris-Mwito5/ginja-ai/internal/db"
	"github.com/Doris-Mwito5/ginja-ai/internal/jwt"
	"github.com/Doris-Mwito5/ginja-ai/internal/services"
	"github.com/gin-gonic/gin"
	"time"
)

const DefaultTokenDuration = 24 * time.Hour

func AddEndpoints(
    public *gin.RouterGroup,    
    protected *gin.RouterGroup, 
    dB db.DB,
    userService services.UserService,
    jwtMaker jwt.JWTToken,
    tokenDuration time.Duration,
) {
    if tokenDuration <= 0 {
        tokenDuration = DefaultTokenDuration
    }

    // Public Endpoints
    public.POST("/register", register(dB, userService))
    public.POST("/login", login(dB, userService, jwtMaker, tokenDuration))

    // Protected Endpoints
    
    protected.GET("/:id", getUserByID(dB, userService))
    protected.GET("/username/:username", getUserByUsername(dB, userService))
}