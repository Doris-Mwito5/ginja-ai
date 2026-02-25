package claims

import (
	"github.com/Doris-Mwito5/ginja-ai/internal/db"
	"github.com/Doris-Mwito5/ginja-ai/internal/services"
	"github.com/gin-gonic/gin"
)

func AddEndpoints(
	r *gin.RouterGroup,
	dB db.DB,
	claimService services.ClaimService,
) {
	r.POST("/claims", createClaim(dB, claimService))
	r.GET("/claims/:id", getClaim(dB, claimService))
	r.GET("/claims/member/:memberID", listClaims(dB, claimService))
}
