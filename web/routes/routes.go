package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/Doris-Mwito5/ginja-ai/internal/configs"
	"github.com/Doris-Mwito5/ginja-ai/internal/db"
	"github.com/Doris-Mwito5/ginja-ai/internal/domain"
	"github.com/Doris-Mwito5/ginja-ai/internal/jwt"
	middleware "github.com/Doris-Mwito5/ginja-ai/internal/middleware"
	"github.com/Doris-Mwito5/ginja-ai/internal/services"
	"github.com/Doris-Mwito5/ginja-ai/web/handlers/claims"
	"github.com/Doris-Mwito5/ginja-ai/web/handlers/members"
	"github.com/Doris-Mwito5/ginja-ai/web/handlers/procedures"
	"github.com/Doris-Mwito5/ginja-ai/web/handlers/providers"
	"github.com/Doris-Mwito5/ginja-ai/web/handlers/users"
	"github.com/gin-gonic/gin"
)

type AppRouter struct {
	*gin.Engine
}

func BuildRouter(
	dB db.DB,
	domainStore *domain.Store,
) *AppRouter {
	router := gin.Default()

	baseAPIGroup := router.Group("/v1")
	baseAPIGroup.Use(middleware.CORSMiddleware())

	jwtMaker, err := jwt.NewJWTMaker(configs.Config.JWTSecret)
	if err != nil {
		log.Fatalf("Failed to create JWT maker: %v", err)
	}

	// --- Service Instantiation ---
	claimService := services.NewClaimService(domainStore)
	userService := services.NewUserService(domainStore)
	memberService := services.NewMemberService(domainStore)
	procedureService := services.NewProcedureService(domainStore)
	providerService := services.NewProviderService(domainStore)

	// Public group (no auth)
	publicRoutes := baseAPIGroup.Group("")

	// Protected group (JWT required)
	protectedRoutes := baseAPIGroup.Group("")
	protectedRoutes.Use(middleware.AuthMiddleware(jwtMaker))

	users.AddEndpoints(publicRoutes, protectedRoutes, dB, userService, jwtMaker, 24*time.Hour)

	claims.AddEndpoints(protectedRoutes, dB, claimService)

	members.AddEndpoints(protectedRoutes, dB, memberService)
	procedures.AddEndpoints(protectedRoutes, dB, procedureService)
	providers.AddEndpoints(protectedRoutes, dB, providerService)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error_message": "Endpoint not found"})
	})

	return &AppRouter{router}
}
