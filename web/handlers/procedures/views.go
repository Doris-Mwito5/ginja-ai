package procedures

import (
	"net/http"

	"github.com/Doris-Mwito5/ginja-ai/internal/apperr"
	"github.com/Doris-Mwito5/ginja-ai/internal/db"
	"github.com/Doris-Mwito5/ginja-ai/internal/dtos"
	"github.com/Doris-Mwito5/ginja-ai/internal/services"
	"github.com/Doris-Mwito5/ginja-ai/internal/utils"
	"github.com/gin-gonic/gin"
)

func createProcedure(
	dB db.DB,
	procedureService services.ProcedureService,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req dtos.Procedure
		err := c.BindJSON(&req)
		if err != nil {
			utils.HandleError(c, apperr.NewErrorWithType(err, apperr.BadRequest))
			return
		}

		procedure, err := procedureService.CreateProcedure(c.Request.Context(), dB, &req)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusCreated, procedure)
	}
}