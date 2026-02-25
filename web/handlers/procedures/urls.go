package procedures

import (
	"github.com/Doris-Mwito5/ginja-ai/internal/db"
	"github.com/Doris-Mwito5/ginja-ai/internal/services"
	"github.com/gin-gonic/gin"
)

func AddEndpoints(
	r *gin.RouterGroup,
	dB db.DB,
	procedureService services.ProcedureService,
) {
	r.POST("/procedures", createProcedure(dB, procedureService))
}